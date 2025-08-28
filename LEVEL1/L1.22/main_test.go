package calc

import (
	"bytes"
	"math/big"
	"os"
	"strings"
	"testing"
)

// bigOpString вычисляет результат операции a (op) b через math/big
// и возвращает строковое представление. Используем в тестах как "эталон".
func bigOpString(a, b int, op string) string {
	A := big.NewInt(int64(a))
	B := big.NewInt(int64(b))
	switch op {
	case "+":
		return big.NewInt(0).Add(A, B).String()
	case "-":
		return big.NewInt(0).Sub(A, B).String()
	case "*":
		return big.NewInt(0).Mul(A, B).String()
	case "/":
		return big.NewInt(0).Div(A, B).String()
	default:
		return ""
	}
}

func TestAdd(t *testing.T) {
	maxInt := int(^uint(0) >> 1)

	tests := []struct {
		name string
		a, b int
	}{
		{"small", 1, 2},
		{"negative", -5, 3},
		{"overflow_max_plus_one", maxInt, 1},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := add(tc.a, tc.b)
			want := bigOpString(tc.a, tc.b, "+")
			if got != want {
				t.Fatalf("add(%d,%d) = %q; want %q", tc.a, tc.b, got, want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	maxInt := int(^uint(0) >> 1)

	tests := []struct {
		name string
		a, b int
	}{
		{"small", 10, 3},
		{"negative_result", 3, 10},
		{"underflow_min_minus_one", -maxInt - 1, 1}, // minInt - 1 -> big
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := subtract(tc.a, tc.b)
			want := bigOpString(tc.a, tc.b, "-")
			if got != want {
				t.Fatalf("subtract(%d,%d) = %q; want %q", tc.a, tc.b, got, want)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	maxInt := int(^uint(0) >> 1)

	tests := []struct {
		name string
		a, b int
	}{
		{"zero", 0, 12345},
		{"small", 7, 6},
		{"overflow_max_times_two", maxInt, 2},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := multiply(tc.a, tc.b)
			want := bigOpString(tc.a, tc.b, "*")
			if got != want {
				t.Fatalf("multiply(%d,%d) = %q; want %q", tc.a, tc.b, got, want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1

	tests := []struct {
		name string
		a, b int
	}{
		{"normal", 10, 3},
		{"divide_by_zero", 10, 0},
		{"min_div_minus_one_overflow", minInt, -1},
		{"negative_division", -9, 3},
		{"exact_division", 12, 3},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := divide(tc.a, tc.b)
			var want string
			if tc.b == 0 {
				want = "cannot divide by 0"
			} else {
				want = bigOpString(tc.a, tc.b, "/")
			}
			if got != want {
				t.Fatalf("divide(%d,%d) = %q; want %q", tc.a, tc.b, got, want)
			}
		})
	}
}

func TestInput(t *testing.T) {
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	t.Run("valid_two_numbers", func(t *testing.T) {
		a, b, err := input(bytes.NewBufferString("42 99\n"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if a != 42 || b != 99 {
			t.Fatalf("input() = %d, %d; want 42, 99", a, b)
		}
	})

	t.Run("too_many_values", func(t *testing.T) {
		_, _, err := input(bytes.NewBufferString("1 2 3\n"))
		if err == nil || !strings.Contains(err.Error(), "expected 2 numbers") {
			t.Fatalf("expected error about expected 2 numbers, got: %v", err)
		}
	})

	t.Run("non_integer", func(t *testing.T) {
		_, _, err := input(bytes.NewBufferString("foo 2\n"))
		if err == nil || !strings.Contains(err.Error(), "Failed to convert input data #1") {
			t.Fatalf("expected conversion error for first value, got: %v", err)
		}
	})
}
