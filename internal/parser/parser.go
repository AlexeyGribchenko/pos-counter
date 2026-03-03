package parser

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	InputFile  string
	OutputFile string
	Help       bool
}

func ParseFlags() (*Config, error) {

	config := &Config{}

	flag.StringVar(&config.InputFile, "input", "", "Input `file path` (required)")
	flag.StringVar(&config.InputFile, "i", "", "Input `file path` (required)")

	flag.StringVar(&config.OutputFile, "output", "", "Output `file path` (optional)")
	flag.StringVar(&config.OutputFile, "o", "", "Output `file path` (optional)")

	flag.BoolVar(&config.Help, "help", false, "Show help")
	flag.BoolVar(&config.Help, "h", false, "Show help")

	flag.Parse()

	if config.Help {
		printUsage()
	}

	if config.InputFile == "" {
		return nil, fmt.Errorf("input file is required")
	}

	return config, nil
}

func printUsage() {
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
}
