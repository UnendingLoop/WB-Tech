package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunScriptFile(t *testing.T) {
	tmpDir := t.TempDir()
	scriptFile := filepath.Join(tmpDir, "script.txt")
	outFile := filepath.Join(tmpDir, "out.txt")

	// Простой скрипт: echo и редирект
	script := `
# Это комментарий
echo Hello > ` + outFile + `
echo World >> ` + outFile + `
`
	if err := os.WriteFile(scriptFile, []byte(script), 0o644); err != nil {
		t.Fatalf("failed to write script: %v", err)
	}

	err := runScriptFile(scriptFile)
	if err != nil {
		t.Fatalf("runScriptFile failed: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	content := strings.TrimSpace(string(data))
	if !strings.Contains(content, "Hello") || !strings.Contains(content, "World") {
		t.Errorf("output file content incorrect: %q", content)
	}
}
