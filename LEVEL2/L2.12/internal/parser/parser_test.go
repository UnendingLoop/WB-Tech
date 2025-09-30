package parser

import (
	"os"
	"testing"
)

func TestInitSearchParam_BasicFlags(t *testing.T) {
	cases := []struct {
		name        string
		args        []string
		wantAfter   int
		wantBefore  int
		wantCircle  int
		wantCount   bool
		wantIgnore  bool
		wantInvert  bool
		wantExact   bool
		wantEnum    bool
		wantPattern string
		wantSources []string
	}{
		{
			name:        "Basic run with pattern and without flags",
			args:        []string{"grepClone", "hello"},
			wantPattern: "hello",
		},
		{
			name:        "A/B/C specified, A/B being padded by C",
			args:        []string{"grepClone", "-A", "2", "-B", "3", "-C", "4", "foo"},
			wantAfter:   4,
			wantBefore:  4,
			wantCircle:  4,
			wantPattern: "foo",
		},
		{
			name:        "Testing flags c, i, v, F, n correctness",
			args:        []string{"grepClone", "-c", "-i", "-v", "-F", "-n", "bar"},
			wantCount:   true,
			wantIgnore:  true,
			wantInvert:  true,
			wantExact:   true,
			wantEnum:    true,
			wantPattern: "bar",
		},
		{
			name:        "C-flag only - padding A and B",
			args:        []string{"grepClone", "-C", "4", "foo"},
			wantAfter:   4,
			wantBefore:  4,
			wantCircle:  4,
			wantPattern: "foo",
		},
		{
			name:        "C-flag between A and B - padding A ",
			args:        []string{"grepClone", "-A", "7", "-C", "3", "-B", "7", "foo"},
			wantAfter:   3,
			wantBefore:  7,
			wantCircle:  3,
			wantPattern: "foo",
		},
		{
			name:        "C-flag between B and A - padding B ",
			args:        []string{"grepClone", "-B", "7", "-C", "3", "-A", "7", "foo"},
			wantAfter:   7,
			wantBefore:  3,
			wantCircle:  3,
			wantPattern: "foo",
		},
		{
			name:        "Pattern + several sources - PrintFileName must become active",
			args:        []string{"grepClone", "zzz", "file1.txt", "file2.txt"},
			wantPattern: "zzz",
			wantSources: []string{"file1.txt", "file2.txt"},
		},
	}

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctxPriority = map[string]int{}
			os.Args = tc.args

			err := InitSearchParam()
			if err != nil {
				t.Fatalf("Fetched error from InitSearchParam(): %v", err)
			}

			if SP.Pattern != tc.wantPattern {
				t.Errorf("Pattern actual = %q; expected %q", SP.Pattern, tc.wantPattern)
			}
			if SP.CtxAfter != tc.wantAfter {
				t.Errorf("CtxAfter actual = %d; expected %d", SP.CtxAfter, tc.wantAfter)
			}
			if SP.CtxBefore != tc.wantBefore {
				t.Errorf("CtxBefore actual = %d; expected %d", SP.CtxBefore, tc.wantBefore)
			}
			if SP.CtxCircle != tc.wantCircle {
				t.Errorf("CtxCircle actual = %d; expected %d", SP.CtxCircle, tc.wantCircle)
			}
			if SP.CountFound != tc.wantCount {
				t.Errorf("CountFound actual = %v; expected %v", SP.CountFound, tc.wantCount)
			}
			if SP.IgnoreCase != tc.wantIgnore {
				t.Errorf("IgnoreCase actual = %v; expected %v", SP.IgnoreCase, tc.wantIgnore)
			}
			if SP.InvertResult != tc.wantInvert {
				t.Errorf("InvertResult actual = %v; expected %v", SP.InvertResult, tc.wantInvert)
			}
			if SP.ExactMatch != tc.wantExact {
				t.Errorf("ExactMatch actual = %v; expected %v", SP.ExactMatch, tc.wantExact)
			}
			if SP.EnumLine != tc.wantEnum {
				t.Errorf("EnumLine actual = %v; expected %v", SP.EnumLine, tc.wantEnum)
			}

			if len(tc.wantSources) > 0 {
				if len(SP.Source) != len(tc.wantSources) {
					t.Fatalf("Source length actual = %d; expected %d", len(SP.Source), len(tc.wantSources))
				}
				for i := range SP.Source {
					if SP.Source[i] != tc.wantSources[i] {
						t.Errorf("Source[%d] actual = %q; expected %q", i, SP.Source[i], tc.wantSources[i])
					}
				}
			}
		})
	}
}

func TestInitSearchParam_Errors(t *testing.T) {
	t.Run("No pattern - must show error", func(t *testing.T) {
		oldArgs := os.Args
		os.Args = []string{"grepClone"}
		defer func() { os.Args = oldArgs }()

		err := InitSearchParam()
		if err == nil {
			t.Fatal("Expected error due to flags and args absence, but got err == nil")
		}
	})

	t.Run("Invalid regexp must show error", func(t *testing.T) {
		oldArgs := os.Args
		os.Args = []string{"grepClone", "("} // незакрытая скобка
		defer func() { os.Args = oldArgs }()

		err := InitSearchParam()
		if err == nil {
			t.Fatal("Expected error whille compiling the regexp, but got err == nil")
		}
	})
}

func TestInitSearchParam_IgnoreCaseRegexp(t *testing.T) {
	oldArgs := os.Args
	os.Args = []string{"grepClone", "-i", "foo"}
	defer func() { os.Args = oldArgs }()

	err := InitSearchParam()
	if err != nil {
		t.Fatalf("InitSearchParam() returned error: %v", err)
	}

	if SP.RegexpPattern == nil {
		t.Fatalf("RegexpPattern must be initialized")
	}

	// Проверим, что установлен (?i)
	if SP.RegexpPattern.String() != "(?i)foo" {
		t.Errorf("RegexpPattern.String() actual = %q; expected %q", SP.RegexpPattern.String(), "(?i)foo")
	}

	// Проверим, что реально матчится без учёта регистра
	if !SP.RegexpPattern.MatchString("FoO") {
		t.Errorf("RegexpPattern must match the string 'FoO'")
	}
}
