package fileops

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestFileExists(t *testing.T) {
	// Создаем временный файл
	tmpFile, err := os.CreateTemp("", "testfile*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	tests := []struct {
		name     string
		filename string
		setup    func() error
		cleanup  func()
		expected bool
	}{
		{
			name:     "Existing file",
			filename: tmpFile.Name(),
			expected: true,
		},
		{
			name:     "Non-existing file",
			filename: "nonexistent123456.txt",
			expected: false,
		},
		{
			name:     "Empty filename",
			filename: "",
			expected: false,
		},
		{
			name:     "Directory instead of file",
			filename: ".",
			expected: false,
		},
		{
			name:     "File with special characters",
			filename: "test@#$%.txt",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FileExists(tt.filename)
			if result != tt.expected {
				t.Errorf("FileExists(%q) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	// Создаем временный файл с тестовым содержимым
	tmpFile, err := os.CreateTemp("", "testread*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	testContent := "Hello, World!\nThis is a test file.\nLine 3"
	_, err = tmpFile.WriteString(testContent)
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	tests := []struct {
		name        string
		filename    string
		setup       func() error
		cleanup     func()
		expected    string
		expectError bool
	}{
		{
			name:        "Read existing file",
			filename:    tmpFile.Name(),
			expected:    testContent,
			expectError: false,
		},
		{
			name:        "Read non-existing file",
			filename:    "nonexistent123456.txt",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Read empty filename",
			filename:    "",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Read directory",
			filename:    ".",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := ReadFile(tt.filename)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if content != tt.expected {
					t.Errorf("ReadFile() = %q, want %q", content, tt.expected)
				}
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	// Создаем временную директорию для тестов
	tmpDir, err := os.MkdirTemp("", "testwrite")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name        string
		filename    string
		content     string
		setup       func() error
		cleanup     func()
		expectError bool
	}{
		{
			name:        "Write to new file",
			filename:    filepath.Join(tmpDir, "newfile.txt"),
			content:     "Test content",
			expectError: false,
		},
		{
			name:     "Write to existing file",
			filename: filepath.Join(tmpDir, "existing.txt"),
			content:  "Updated content",
			setup: func() error {
				return os.WriteFile(filepath.Join(tmpDir, "existing.txt"), []byte("Original content"), 0644)
			},
			expectError: false,
		},
		{
			name:        "Write empty content",
			filename:    filepath.Join(tmpDir, "empty.txt"),
			content:     "",
			expectError: false,
		},
		{
			name:        "Write to invalid path",
			filename:    "/invalid/path/that/does/not/exist/file.txt",
			content:     "Test content",
			expectError: true,
		},
		{
			name:        "Write with empty filename",
			filename:    "",
			content:     "Test content",
			expectError: true,
		},
		{
			name:        "Write large content",
			filename:    filepath.Join(tmpDir, "large.txt"),
			content:     string(make([]byte, 10000)), // 10KB of null bytes
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Выполняем setup если есть
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Выполняем тест
			err := WriteFile(tt.filename, tt.content)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Проверяем, что файл был создан и содержит правильное содержимое
				if !tt.expectError && tt.filename != "" {
					content, readErr := os.ReadFile(tt.filename)
					if readErr != nil {
						t.Errorf("Failed to read written file: %v", readErr)
					}
					if string(content) != tt.content {
						t.Errorf("File content = %q, want %q", string(content), tt.content)
					}
				}
			}
		})
	}
}

func TestReadWriteIntegration(t *testing.T) {
	// Интеграционный тест: записываем файл, затем читаем его
	tmpFile, err := os.CreateTemp("", "integration*.txt")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	originalContent := "Integration test content\nLine 2\nLine 3"

	// Записываем
	err = WriteFile(tmpFile.Name(), originalContent)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Читаем
	readContent, err := ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if readContent != originalContent {
		t.Errorf("Content mismatch:\ngot: %q\nwant: %q", readContent, originalContent)
	}
}

func BenchmarkReadFile(b *testing.B) {
	// Создаем тестовый файл
	tmpFile, err := os.CreateTemp("", "benchread*.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := strings.Repeat("Benchmark line\n", 1000)
	_, err = tmpFile.WriteString(content)
	if err != nil {
		b.Fatal(err)
	}
	tmpFile.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ReadFile(tmpFile.Name())
	}
}

func BenchmarkWriteFile(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "benchwrite")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	content := strings.Repeat("Benchmark line\n", 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filename := filepath.Join(tmpDir, "bench_"+strconv.Itoa(i)+".txt")
		WriteFile(filename, content)
	}
}
