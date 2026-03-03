package main

import (
	"fmt"
	"log"

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

	fmt.Printf("Input file: %s", cfg.InputFile)
	if cfg.OutputFile != "" {
		fmt.Printf("Output file: %s", cfg.OutputFile)
	}

}
