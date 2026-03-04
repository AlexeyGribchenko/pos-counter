package main

import (
	"fmt"
	"log"

	"github.com/AlexeyGribchenko/pos-counter/internal/counter"
	"github.com/AlexeyGribchenko/pos-counter/internal/fileops"
	"github.com/AlexeyGribchenko/pos-counter/internal/parser"
)

func main() {
	cfg, err := parser.ParseFlags()
	if err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	if cfg == nil {
		return
	}

	if !fileops.FileExists(cfg.InputFile) {
		log.Fatalf("Input file does not exist: %s", cfg.InputFile)
	}

	text, err := fileops.ReadFile(cfg.InputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	posCounter := counter.NewPOSConuter()
	result := posCounter.Count(text)

	if cfg.OutputFile != "" {
		err := fileops.WriteFile(cfg.OutputFile, result.String())
		if err != nil {
			log.Fatalf("Error writing output file: %v", err)
		}
		fmt.Println("Result written to", cfg.OutputFile)
	} else {
		fmt.Println(result)
	}
}
