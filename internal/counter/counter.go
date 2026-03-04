package counter

import (
	"regexp"
	"strings"

	"github.com/AlexeyGribchenko/pos-counter/internal/models"
)

type POSCounter struct {
	adjectiveRegex *regexp.Regexp
	adverbRegex    *regexp.Regexp
	verbRegex      *regexp.Regexp
	wordRegex      *regexp.Regexp
}

func NewPOSConuter() *POSCounter {

	adjectivePattern := `(?i)\b\w+(?:y|ous|ive|ful|less|ic|al|able|ible|ish|like|ly|an|ese|que)\b`
	adverbPattern := `(?i)\b\w+ly\b`
	verbPattern := `(?i)\b\w+(?:ing|ed|en|s|es|ies|ate|ize|ise|ify|en)\b`
	wordPattern := `\b[a-zA-Z]+\b`

	return &POSCounter{
		adjectiveRegex: regexp.MustCompile(adjectivePattern),
		adverbRegex:    regexp.MustCompile(adverbPattern),
		verbRegex:      regexp.MustCompile(verbPattern),
		wordRegex:      regexp.MustCompile(wordPattern),
	}
}

func (p *POSCounter) Count(text string) models.POSresult {

	text = strings.ToLower(text)

	return models.POSresult{
		Adjectives: len(p.adjectiveRegex.FindAllString(text, -1)),
		Adverbs:    len(p.adverbRegex.FindAllString(text, -1)),
		Verbs:      len(p.verbRegex.FindAllString(text, -1)),
		Words:      len(p.wordRegex.FindAllString(text, -1)),
	}
}
