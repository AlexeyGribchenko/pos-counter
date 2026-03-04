# Part of Speech Counter

A Go application that counts adjectives, adverbs, and verbs in English text files.

## Features
- Counts adjectives, adverbs, and verbs
- Detailed word lists by category
- Input/output file support
- Unit tests

## Limitations and Known Issues

### Regular Expression Limitations
The application uses pattern-based recognition rather than full grammatical parsing. This approach has inherent limitations:

- **Overlapping patterns**: Words ending in "ly" can be either adverbs (quickly, beautifully) or adjectives (lovely, friendly). The current implementation may misclassify such words.
- **False positives**: Words that match suffix patterns but belong to other parts of speech might be incorrectly counted. For example, "supply" matches the adjective pattern but is primarily a verb/noun.
- **Context independence**: The counter analyzes words in isolation, ignoring grammatical context. The same word might be different parts of speech depending on usage (e.g., "book" can be noun or verb).

### Example of Misclassification
Input: "The friendly girl sings beautifully"
- "friendly" (adjective) might be counted as an adverb due to "-ly" suffix
- "beautifully" (adverb) correctly identified but highlights pattern overlap

## Quick Start
```bash
go build -o pos-counter ./cmd/pos-counter/main.go
./pos-counter -i ./testdata/input.txt