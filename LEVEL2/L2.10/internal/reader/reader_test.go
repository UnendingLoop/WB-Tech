package reader_test

import (
	"os"
	"strings"
	"testing"

	"sortClone/internal/reader"
)

func TestReadInput_SmallFile(t *testing.T) {
	// создаём временный файл с несколькими строками
	tmpFile := t.TempDir() + "/small.txt"
	content := "one\ntwo\nthree\n"
	if err := os.WriteFile(tmpFile, []byte(content), 0o644); err != nil {
		t.Fatalf("cannot create temp file: %v", err)
	}

	lines, tmpFiles, err := reader.ReadInput([]string{tmpFile})
	if err != nil {
		t.Fatalf("ReadInput failed: %v", err)
	}

	if tmpFiles != nil {
		t.Errorf("expected no tmp files, got %v", tmpFiles)
	}

	expected := []string{"one", "two", "three"}
	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("expected line %d = %q, got %q", i, expected[i], line)
		}
	}
}

func TestReadInput_LargeFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/large.txt"

	// создаём "большой" файл, который больше maxFileSize
	var b strings.Builder
	for i := 0; i < 1024*1024; i++ { // ~1 млн строк
		b.WriteString("lineфывапролджэйцукенгшщзхъячсмитьбюйцукенгшщзхъфывапролджэячсмитьбюйцукенгшщзхъфывапролджэячсмитьбювап\n")
	}

	if err := os.WriteFile(tmpFile, []byte(b.String()), 0o644); err != nil {
		t.Fatalf("cannot create large temp file: %v", err)
	}

	_, tmpFiles, err := reader.ReadInput([]string{tmpFile})
	if err != nil {
		t.Fatalf("ReadInput failed: %v", err)
	}

	if len(tmpFiles) == 0 {
		t.Fatalf("expected tmp files to be created, got none")
	}

	// проверим, что временные файлы существуют и не пусты
	for _, f := range tmpFiles {
		info, err := os.Stat(f)
		if err != nil {
			t.Errorf("tmp file %s does not exist: %v", f, err)
		}
		if info.Size() == 0 {
			t.Errorf("tmp file %s is empty", f)
		}
	}

	os.RemoveAll("tmp")
}
