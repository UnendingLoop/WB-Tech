package reader_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"cutClone/internal/parser"
	"cutClone/internal/reader"
)

func captureStdout(t *testing.T) (*bytes.Buffer, func()) {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() failed: %v", err)
	}

	old := os.Stdout
	os.Stdout = w

	var buf bytes.Buffer
	done := make(chan struct{})

	// горутина, которая читает из pipe в буфер
	go func() {
		_, _ = io.Copy(&buf, r)
		close(done)
	}()

	restore := func() {
		w.Close()
		os.Stdout = old
		<-done
	}

	return &buf, restore
}

func tempFileWithContent(t *testing.T, content string) string {
	t.Helper()

	tmp := filepath.Join(t.TempDir(), "input.txt")
	if err := os.WriteFile(tmp, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	return tmp
}

func TestProcessInput_FromStdin(t *testing.T) {
	parser.CC.Delim = " "
	parser.CC.SepOnly = false
	parser.CC.Fields = "1"

	// попдмена stdin
	input := "one two three\nfoo bar\n"
	r, w, _ := os.Pipe()
	oldStdin := os.Stdin
	os.Stdin = r

	// захват stdout
	buf, restore := captureStdout(t)

	go func() {
		_, _ = w.Write([]byte(input))
		w.Close()
	}()

	err := reader.ProcessInput("")
	if err != nil {
		t.Fatalf("ProcessInput returned error: %v", err)
	}

	restore()
	os.Stdin = oldStdin

	out := buf.String()
	if out == "" {
		t.Errorf("expected some output, got empty string")
	}
}

func TestProcessInput_FromFile(t *testing.T) {
	parser.CC.Delim = " "
	parser.CC.SepOnly = true
	parser.CC.Fields = "1"

	content := "alpha beta\ngamma delta\n"
	file := tempFileWithContent(t, content)

	buf, restore := captureStdout(t)
	defer restore()

	err := reader.ProcessInput(file)
	if err != nil {
		t.Fatalf("ProcessInput returned error: %v", err)
	}

	out := buf.String()
	if out == "" {
		t.Errorf("expected some output, got empty string")
	}
}
