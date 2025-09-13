package utils

import (
	"strconv"
	"testing"

	"sortClone/internal/model"
)

func TestGetKey(t *testing.T) {
	tests := []struct {
		line     string
		column   int
		delim    string
		expected string
	}{
		{"a b c", 1, " ", "a"},
		{"a b c", 2, " ", "b"},
		{"a b c", 3, " ", "c"},
		{"a b c", 4, " ", "a b c"}, // колонка вне диапазона
		{"x,y,z", 2, ",", "y"},
		{"no-delim", 1, ",", "no-delim"},
		{"no-delim", 2, ",", "no-delim"},
	}

	for _, tt := range tests {
		got := GetKey(tt.line, tt.column, tt.delim)
		if got != tt.expected {
			t.Errorf("GetKey(%q, %d, %q) = %q, want %q",
				tt.line, tt.column, tt.delim, got, tt.expected)
		}
	}
}

func TestUniqueLines(t *testing.T) {
	lines := []string{
		"apple 1",
		"banana 2",
		"apple 3", // повторяется по первой колонке
		"banana 4",
		"cherry 5",
	}
	opts := model.Options{Column: 1, Delimeter: " "}
	got := UniqueLines(lines, opts)

	expected := []string{"apple 1", "banana 2", "cherry 5"}
	if len(got) != len(expected) {
		t.Fatalf("UniqueLines() = %v, want %v", got, expected)
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("UniqueLines()[%d] = %q, want %q", i, got[i], expected[i])
		}
	}
}

func TestRemoveTrailBlanks(t *testing.T) {
	a, b := RemoveTrailBlanks("hello   ", "world  ")
	if a != "hello" || b != "world" {
		t.Errorf("RemoveTrailBlanks() = %q, %q; want %q, %q",
			a, b, "hello", "world")
	}
}

func TestParseHumanSize(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"1024", 1024},
		{"1K", 1024},
		{"2k", 2048}, // регистр неважен
		{"1M", 1024 * 1024},
		{"1.5M", 1.5 * 1024 * 1024},
		{"2G", 2 * 1024 * 1024 * 1024},
		{"", 0},
	}

	for _, tt := range tests {
		got, err := ParseHumanSize(tt.input)
		if err != nil {
			t.Errorf("ParseHumanSize(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if got != tt.expected {
			t.Errorf("ParseHumanSize(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestSmallComparator(t *testing.T) {
	// числовое сравнение
	if !SmallComparator("1", "2", 1, 2, nil, nil) {
		t.Errorf("expected 1 < 2")
	}
	if SmallComparator("2", "1", 2, 1, nil, nil) {
		t.Errorf("expected 2 > 1")
	}

	// оба с ошибкой — сравнение строк
	if !SmallComparator("apple", "banana", 0, 0, assertErr(), assertErr()) {
		t.Errorf("expected apple < banana by string comparison")
	}

	// только первая с ошибкой
	if SmallComparator("apple", "123", 0, 123, assertErr(), nil) {
		t.Errorf("expected string with error to be greater")
	}

	// только вторая с ошибкой
	if !SmallComparator("123", "apple", 123, 0, nil, assertErr()) {
		t.Errorf("expected valid number to be less than error string")
	}
}

// вспомогательная функция для имитации ошибки
func assertErr() error {
	return &strconv.NumError{Func: "Atoi", Num: "x", Err: strconv.ErrSyntax}
}
