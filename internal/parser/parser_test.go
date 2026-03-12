package parser

import (
	"flag"
	"os"
	"testing"
)

func TestParseFlags(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tests := []struct {
		name        string
		args        []string
		expected    *Config
		expectError bool
	}{
		{
			name: "Input file only with long flag",
			args: []string{"cmd", "-input", "test.txt"},
			expected: &Config{
				InputFile:  "test.txt",
				OutputFile: "",
				Help:       false,
			},
			expectError: false,
		},
		{
			name: "Input file only with short flag",
			args: []string{"cmd", "-i", "test.txt"},
			expected: &Config{
				InputFile:  "test.txt",
				OutputFile: "",
				Help:       false,
			},
			expectError: false,
		},
		{
			name: "Input and output files with long flags",
			args: []string{"cmd", "-input", "in.txt", "-output", "out.txt"},
			expected: &Config{
				InputFile:  "in.txt",
				OutputFile: "out.txt",
				Help:       false,
			},
			expectError: false,
		},
		{
			name: "Input and output files with short flags",
			args: []string{"cmd", "-i", "in.txt", "-o", "out.txt"},
			expected: &Config{
				InputFile:  "in.txt",
				OutputFile: "out.txt",
				Help:       false,
			},
			expectError: false,
		},
		{
			name: "Mixed flags",
			args: []string{"cmd", "-input", "in.txt", "-o", "out.txt"},
			expected: &Config{
				InputFile:  "in.txt",
				OutputFile: "out.txt",
				Help:       false,
			},
			expectError: false,
		},
		{
			name: "Help flag long",
			args: []string{"cmd", "-help", "-i", "input.txt"},
			expected: &Config{
				InputFile:  "input.txt",
				OutputFile: "",
				Help:       true,
			},
			expectError: false,
		},
		{
			name: "Help flag short",
			args: []string{"cmd", "-h", "-i", "input.txt"},
			expected: &Config{
				InputFile:  "input.txt",
				OutputFile: "",
				Help:       true,
			},
			expectError: false,
		},
		{
			name:        "No input file",
			args:        []string{"cmd"},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Input file with spaces",
			args: []string{"cmd", "-input", "my test file.txt"},
			expected: &Config{
				InputFile:  "my test file.txt",
				OutputFile: "",
				Help:       false,
			},
			expectError: false,
		},
		{
			name: "All flags with help",
			args: []string{"cmd", "-i", "in.txt", "-o", "out.txt", "-h"},
			expected: &Config{
				InputFile:  "in.txt",
				OutputFile: "out.txt",
				Help:       true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = tt.args

			config, err := ParseFlags()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.expected == nil {
				if config != nil {
					t.Errorf("Expected nil config, got %+v", config)
				}
				return
			}

			if config == nil {
				t.Fatal("Expected config, got nil")
			}

			if config.InputFile != tt.expected.InputFile {
				t.Errorf("InputFile = %q, want %q", config.InputFile, tt.expected.InputFile)
			}
			if config.OutputFile != tt.expected.OutputFile {
				t.Errorf("OutputFile = %q, want %q", config.OutputFile, tt.expected.OutputFile)
			}
			if config.Help != tt.expected.Help {
				t.Errorf("Help = %v, want %v", config.Help, tt.expected.Help)
			}
		})
	}
}

func TestParseFlags_HelpDoesNotReturnError(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-h", "-i", "input.txt"}

	config, err := ParseFlags()

	if err != nil {
		t.Errorf("Expected no error for help flag, got %v", err)
	}

	if config == nil {
		t.Fatal("Expected config for help flag, got nil")
	}

	if !config.Help {
		t.Error("Expected Help=true, got false")
	}
}

func TestParseFlags_EdgeCases(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tests := []struct {
		name        string
		args        []string
		checkConfig func(*testing.T, *Config)
		expectError bool
	}{
		{
			name: "Input file with equals sign",
			args: []string{"cmd", "-input=test.txt"},
			checkConfig: func(t *testing.T, c *Config) {
				if c.InputFile != "test.txt" {
					t.Errorf("InputFile = %q, want %q", c.InputFile, "test.txt")
				}
			},
			expectError: false,
		},
		{
			name: "Multiple input flags (last wins)",
			args: []string{"cmd", "-i", "first.txt", "-input", "second.txt"},
			checkConfig: func(t *testing.T, c *Config) {
				if c.InputFile != "second.txt" {
					t.Errorf("InputFile = %q, want %q", c.InputFile, "second.txt")
				}
			},
			expectError: false,
		},
		{
			name: "Input file with special characters",
			args: []string{"cmd", "-i", "file-with-dashes_and_underscores.txt"},
			checkConfig: func(t *testing.T, c *Config) {
				if c.InputFile != "file-with-dashes_and_underscores.txt" {
					t.Errorf("InputFile = %q, want %q", c.InputFile, "file-with-dashes_and_underscores.txt")
				}
			},
			expectError: false,
		},
		// unknown flag just fails the test
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = tt.args

			config, err := ParseFlags()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if config == nil {
				t.Fatal("Expected config, got nil")
			}

			if tt.checkConfig != nil {
				tt.checkConfig(t, config)
			}
		})
	}
}

func BenchmarkParseFlags(b *testing.B) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	for i := 0; i < b.N; i++ {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		os.Args = []string{"cmd", "-i", "test.txt", "-o", "output.txt"}
		ParseFlags()
	}
}
