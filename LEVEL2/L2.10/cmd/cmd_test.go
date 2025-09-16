package cmd

import (
	"bytes"
	"os"
	"testing"

	"sortClone/internal/model"
)

// Тест сортировки простых строк в памяти
func TestExecuteSort_LinesMemory(t *testing.T) {
	// Подготовка входных данных
	lines := []string{"b", "c", "a"}
	ReadInputFunc = func(args []string) ([]string, []string, error) {
		return lines, nil, nil
	}

	// Перенаправляем вывод в буфер
	buf := &bytes.Buffer{}
	OutputDST = buf
	model.OptsContainer.WriteToFile = "" // Чтобы не пытался писать в файл
	model.OptsContainer.Unique = false

	// Выполнение сортировки
	if err := ExecuteSort(nil, nil); err != nil {
		t.Fatalf("executeSort returned error: %v", err)
	}

	// Проверка результата
	got := buf.String()
	want := "a\nb\nc\n"
	if got != want {
		t.Errorf("unexpected output: got %q, want %q", got, want)
	}
}

// Тест сортировки с уникальными строками
func TestExecuteSort_Unique(t *testing.T) {
	lines := []string{"a", "b", "a", "c"}
	ReadInputFunc = func(args []string) ([]string, []string, error) {
		return lines, nil, nil
	}

	buf := &bytes.Buffer{}
	OutputDST = buf
	model.OptsContainer.WriteToFile = ""
	model.OptsContainer.Unique = true

	if err := ExecuteSort(nil, nil); err != nil {
		t.Fatalf("executeSort returned error: %v", err)
	}

	got := buf.String()
	want := "a\nb\nc\n"
	if got != want {
		t.Errorf("unexpected output: got %q, want %q", got, want)
	}
}

// Тест сортировки временных файлов
func TestExecuteSort_TmpFiles(t *testing.T) {
	// Создание временной директории и файлов
	err := os.MkdirAll("testdata", 0o755)
	if err != nil {
		t.Errorf("Failed to create folder: %v", err)
	}
	defer os.RemoveAll("testdata")

	file1, err := os.CreateTemp("", "tmp1-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file1.Name())

	file2, err := os.CreateTemp("", "tmp2-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file2.Name())

	err = os.WriteFile(file1.Name(), []byte("C\nB\nA\n"), 0o644)
	if err != nil {
		t.Errorf("Failed to write to file1 %q: %v", file1.Name(), err)
	}
	err = os.WriteFile(file2.Name(), []byte("F\nE\nD\n"), 0o644)
	if err != nil {
		t.Errorf("Failed to write to file2 %q: %v", file2.Name(), err)
	}

	err = file1.Close()
	if err != nil {
		t.Errorf("Failed to close file1: %v", err)
	}
	err = file2.Close()
	if err != nil {
		t.Errorf("Failed to close file2: %v", err)
	}
	tmpFiles := []string{file1.Name(), file2.Name()}
	ReadInputFunc = func(args []string) ([]string, []string, error) {
		return nil, tmpFiles, nil
	}

	buf := &bytes.Buffer{}
	OutputDST = buf
	model.OptsContainer.WriteToFile = ""
	model.OptsContainer.Unique = false

	if err := ExecuteSort(nil, nil); err != nil {
		t.Fatalf("executeSort returned error: %v", err)
	}

	got := buf.String()
	want := "A\nB\nC\nD\nE\nF\n"
	if got != want {
		t.Errorf("unexpected output: got %q, want %q", got, want)
	}
}
