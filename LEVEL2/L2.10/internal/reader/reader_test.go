package reader_test

import (
	"bufio"
	"io"
	"os"
	"testing"

	"sortClone/internal/reader"
)

func TestReadInput_File(t *testing.T) {
	// создаём временный файл
	content := "line1\nline2\nline3"
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// вызываем ReadInput
	lines, err := reader.ReadInput([]string{tmpFile.Name()})
	if err != nil {
		t.Fatalf("ReadInput returned error: %v", err)
	}

	expected := []string{"line1", "line2", "line3"}
	if len(lines) != len(expected) {
		t.Fatalf("Expected %d lines, got %d", len(expected), len(lines))
	}
	for i := range lines {
		if lines[i] != expected[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expected[i], lines[i])
		}
	}
}

func ReadInputWithReader(args []string, r io.Reader) ([]string, error) {
	var scanner *bufio.Scanner
	_, err := os.Stat(args[0])
	if len(args) == 1 && err == nil {
		f, err := os.Open(args[0])
		if err != nil {
			return nil, err
		}
		defer f.Close()
		scanner = bufio.NewScanner(f)
	} else {
		scanner = bufio.NewScanner(r)
	}

	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
