package integration

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var testFile1Path, testFile2Path string

func TestStdIn(t *testing.T) {
	cmd := exec.Command("go", "run", "../main.go", "-A", "1", "-i", "0101")

	// пердаём данные в STDIN
	cmd.Stdin = strings.NewReader("0101\n0202\n3030\n0101\n4040")

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	expected := `0101
0202
--
0101
4040
`
	got := out.String()

	if got != expected {
		t.Errorf("unexpected output for STDIN:\n--- got ---\n%s\n--- want ---\n%s", got, expected)
	}
}

func TestSingleFile(t *testing.T) {
	cmd := exec.Command("go", "run", "../main.go", "-A", "1", "-i", "0101", testFile1Path)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	expected := `0101
0202
`

	got := out.String()

	if got != expected {
		t.Errorf("unexpected output for single file:\n--- got ---\n%s\n--- want ---\n%s", got, expected)
	}
}

func TestDualFile(t *testing.T) {
	cmd := exec.Command("go", "run", "../main.go", "-A", "1", "-i", "0101", testFile1Path, testFile2Path)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	expected := testFile1Path + `: 0101
` + testFile1Path + `: 0202
` + testFile2Path + `: 0101
` + testFile2Path + `: 0202
`

	got := out.String()

	if got != expected {
		t.Errorf("unexpected output for dual file:\n--- got ---\n%s\n--- want ---\n%s", got, expected)
	}
}

func TestMain(m *testing.M) {
	content1 := "0101\n0202\nAAAAAAAAAA\n0404\nAAAAAAAAAA\n0606\n0707\nAAAAAAAAAA\n0909\n1010\n1111\nAAAAAAAAAA\nAAAAAAAAAA\nAAAAAAAAAA"
	content2 := "0101\n0202\nBBBBBBBBBB\n0404\nBBBBBBBBBB\n0606\n0707\nBBBBBBBBBB\n0909\n1010\n1111\nBBBBBBBBBB\nBBBBBBBBBB\nBBBBBBBBBB"
	// генерим тестовыe файлы с мокамии
	var err error
	testFile1Path, err = createTempFile(content1)
	if err != nil {
		log.Printf("Failed to create 1st mock-file: %v", err)
		os.Exit(1)
	}
	testFile2Path, err = createTempFile(content2)
	if err != nil {
		log.Printf("Failed to create 2nd mock-file: %v", err)
		os.Exit(1)
	}

	// запускаем все тесты в пакете
	code := m.Run()

	// чистим хлам после тестов
	if err := os.Remove(testFile1Path); err != nil {
		log.Printf("Failed to remove 1st mock-file: %v", err)
		os.Exit(1)
	}
	if err := os.Remove(testFile2Path); err != nil {
		log.Printf("Failed to remove 2nd mock-file: %v", err)
		os.Exit(1)
	}
	os.Exit(code)
}

func createTempFile(content string) (string, error) {
	f, err := os.CreateTemp("", "grepClone_test_*.txt")
	if err != nil {
		return "", err
	}
	if _, err := f.WriteString(content); err != nil {
		return "", err
	}
	f.Close()
	return f.Name(), nil
}
