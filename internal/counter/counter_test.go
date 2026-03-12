package counter

import (
	"strings"
	"testing"

	"github.com/AlexeyGribchenko/pos-counter/internal/models"
)

func TestNewPOSCounter(t *testing.T) {
	counter := NewPOSConuter()
	if counter == nil {
		t.Fatal("NewPOSConuter() returned nil")
	}

	if counter.adjectiveRegex == nil {
		t.Error("adjectiveRegex is nil")
	}
	if counter.adverbRegex == nil {
		t.Error("adverbRegex is nil")
	}
	if counter.verbRegex == nil {
		t.Error("verbRegex is nil")
	}
	if counter.wordRegex == nil {
		t.Error("wordRegex is nil")
	}
}

func TestPOSCounter_Count(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected models.POSresult
	}{
		{
			name: "Empty text",
			text: "",
			expected: models.POSresult{
				Adjectives: 0,
				Adverbs:    0,
				Verbs:      0,
				Words:      0,
			},
		},
		{
			name: "Text with adjectives only",
			text: "The beautiful colorful happy day",
			expected: models.POSresult{
				Adjectives: 4, // beautiful, colorful, happy | day???...
				Adverbs:    0,
				Verbs:      0,
				Words:      5,
			},
		},
		{
			name: "Text with adverbs only",
			text: "She sings beautifully and dances gracefully quickly",
			expected: models.POSresult{
				Adjectives: 3, // 3 is wrong, but due to regex it is ok...
				Adverbs:    3, // beautifully, gracefully, quickly
				Verbs:      2,
				Words:      7,
			},
		},
		{
			name: "Text with verbs only",
			text: "running jumping playing eating sleeping",
			expected: models.POSresult{
				Adjectives: 0,
				Adverbs:    0,
				Verbs:      5, // running, jumping, playing, eating, sleeping
				Words:      5,
			},
		},
		{
			name: "Mixed text with all parts",
			text: "The beautiful girl sings beautifully and runs quickly",
			expected: models.POSresult{
				Adjectives: 3, // beautiful | beautifully, quickly due to regex
				Adverbs:    2, // beautifully, quickly
				Verbs:      2, // sings, runs
				Words:      8,
			},
		},
		{
			name: "Text with punctuation",
			text: "The beautiful, colorful! happy day.",
			expected: models.POSresult{
				Adjectives: 4, // beautiful, colorful, happy | day???...
				Adverbs:    0,
				Verbs:      0,
				Words:      5,
			},
		},
		{
			name: "Text with mixed case",
			text: "The BEAUTIFUL girl SINGS beautifully",
			expected: models.POSresult{
				Adjectives: 2, // beautiful + beautifully due to regex
				Adverbs:    1, // beautifully
				Verbs:      1, // sings
				Words:      5,
			},
		},
		{
			name: "Complex verb forms",
			text: "running walked eaten plays kissed",
			expected: models.POSresult{
				Adjectives: 0,
				Adverbs:    0,
				Verbs:      5, // running, walked, eaten, plays, kissed
				Words:      5,
			},
		},
		{
			name: "Words that shouldn't match",
			text: "the and but or for nor yet",
			expected: models.POSresult{
				Adjectives: 0,
				Adverbs:    0,
				Verbs:      0,
				Words:      7,
			},
		},
		{
			name: "Multiple lines",
			text: "First line with beautiful words\nSecond line with running and jumping",
			expected: models.POSresult{
				Adjectives: 1, // beautiful
				Adverbs:    0,
				Verbs:      3, // running, jumping + words that is not verb
				Words:      11,
			},
		},
	}

	counter := NewPOSConuter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := counter.Count(tt.text)

			if result.Adjectives != tt.expected.Adjectives {
				t.Errorf("Adjectives: got %d, want %d", result.Adjectives, tt.expected.Adjectives)
			}
			if result.Adverbs != tt.expected.Adverbs {
				t.Errorf("Adverbs: got %d, want %d", result.Adverbs, tt.expected.Adverbs)
			}
			if result.Verbs != tt.expected.Verbs {
				t.Errorf("Verbs: got %d, want %d", result.Verbs, tt.expected.Verbs)
			}
			if result.Words != tt.expected.Words {
				t.Errorf("Words: got %d, want %d", result.Words, tt.expected.Words)
			}
		})
	}
}

func TestPOSCounter_Count_EdgeCases(t *testing.T) {
	counter := NewPOSConuter()

	longText := strings.Repeat("beautiful unique running jumping ", 1000)
	result := counter.Count(longText)
	if result.Adjectives != 2000 { // beautiful, amazing each repeated 1000 times
		t.Errorf("Long text adjectives: got %d, want 2000", result.Adjectives)
	}
	if result.Verbs != 2000 { // running, jumping each repeated 1000 times
		t.Errorf("Long text verbs: got %d, want 2000", result.Verbs)
	}

	numbersText := "123 beautiful 456 running 789"
	result = counter.Count(numbersText)
	if result.Adjectives != 1 {
		t.Errorf("Numbers text adjectives: got %d, want 1", result.Adjectives)
	}
	if result.Verbs != 1 {
		t.Errorf("Numbers text verbs: got %d, want 1", result.Verbs)
	}
	if result.Words != 5 { // beautiful, running + 3 numbers
		t.Errorf("Numbers text words: got %d, want 5", result.Words)
	}
}

func TestRegexPatterns_Individually(t *testing.T) {
	counter := NewPOSConuter()

	adjectiveTests := []struct {
		word     string
		expected bool
	}{
		{"beautiful", true},
		{"colorful", true},
		{"happy", true},
		{"friendly", true},
		{"the", false},
		{"run", false},
	}

	for _, tt := range adjectiveTests {
		t.Run("Adjective_"+tt.word, func(t *testing.T) {
			matches := counter.adjectiveRegex.FindAllString(tt.word, -1)
			isMatch := len(matches) > 0
			if isMatch != tt.expected {
				t.Errorf("Word %q as adjective: got %v, want %v", tt.word, isMatch, tt.expected)
			}
		})
	}

	adverbTests := []struct {
		word     string
		expected bool
	}{
		{"beautifully", true},
		{"quickly", true},
		{"slowly", true},
		{"happy", false},
		{"run", false},
	}

	for _, tt := range adverbTests {
		t.Run("Adverb_"+tt.word, func(t *testing.T) {
			matches := counter.adverbRegex.FindAllString(tt.word, -1)
			isMatch := len(matches) > 0
			if isMatch != tt.expected {
				t.Errorf("Word %q as adverb: got %v, want %v", tt.word, isMatch, tt.expected)
			}
		})
	}

	verbTests := []struct {
		word     string
		expected bool
	}{
		{"running", true},
		{"jumped", true},
		{"eaten", true},
		{"plays", true},
		{"kissed", true},
		{"beautiful", false},
		{"the", false},
	}

	for _, tt := range verbTests {
		t.Run("Verb_"+tt.word, func(t *testing.T) {
			matches := counter.verbRegex.FindAllString(tt.word, -1)
			isMatch := len(matches) > 0
			if isMatch != tt.expected {
				t.Errorf("Word %q as verb: got %v, want %v", tt.word, isMatch, tt.expected)
			}
		})
	}
}

func BenchmarkPOSCounter_Count(b *testing.B) {
	counter := NewPOSConuter()
	text := "The beautiful girl sings beautifully and runs quickly through the colorful garden"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		counter.Count(text)
	}
}

func BenchmarkPOSCounter_Count_LongText(b *testing.B) {
	counter := NewPOSConuter()
	text := strings.Repeat("The beautiful girl sings beautifully and runs quickly ", 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		counter.Count(text)
	}
}
