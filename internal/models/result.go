package models

import "fmt"

type POSresult struct {
	Adjectives int
	Adverbs    int
	Verbs      int
	Words      int
}

func (p POSresult) String() string {
	return fmt.Sprintf("POS result:"+
		"\n  Adjectives:\t%d"+
		"\n  Adverbs:\t%d"+
		"\n  Verbs:\t%d"+
		"\n  Total words:\t%d",
		p.Adjectives, p.Adverbs, p.Verbs, p.Words)
}
