# Part of Speech Counter

A Go application that counts adjectives, adverbs, and verbs in English text files.

## Features
- Counts adjectives, adverbs, and verbs
- Detailed word lists by category
- Input/output file support
- Unit tests

## Quick Start
```bash
go build -o pos-counter ./cmd/pos-counter/main.go
./pos-counter -i ./testdata/input.txt