package parser

import (
	"os"
	"testing"
)

func TestInitSearchParam_ValidFlags(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	cases := []struct {
		name        string
		args        []string
		expFields   string
		expDelim    string
		expSepOnly  bool
		expSource   []string
		expectError bool
	}{
		{
			name:       "basic fields",
			args:       []string{"cmd", "-f", "1,2,5", "file.txt"},
			expFields:  "1,2,5",
			expDelim:   "\t",
			expSepOnly: false,
			expSource:  []string{"file.txt"},
		},
		{
			name:       "custom delimiter",
			args:       []string{"cmd", "-f", "1-3", "-d", ",", "file.txt"},
			expFields:  "1-3",
			expDelim:   ",",
			expSepOnly: false,
			expSource:  []string{"file.txt"},
		},
		{
			name:       "suppress lines without delimiter",
			args:       []string{"cmd", "-f", "1-2", "-s", "file.txt"},
			expFields:  "1-2",
			expDelim:   "\t",
			expSepOnly: true,
			expSource:  []string{"file.txt"},
		},
		{
			name:        "invalid column range",
			args:        []string{"cmd", "-f", "abc", "file.txt"},
			expectError: true,
		},
		{
			name:        "missing input",
			args:        []string{"cmd", "-f", "1-3"},
			expectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			os.Args = tc.args

			err := InitSearchParam()
			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if CC.Fields != tc.expFields {
				t.Errorf("Fields = %q; want %q", CC.Fields, tc.expFields)
			}
			if CC.Delim != tc.expDelim {
				t.Errorf("Delim = %q; want %q", CC.Delim, tc.expDelim)
			}
			if CC.SepOnly != tc.expSepOnly {
				t.Errorf("SepOnly = %v; want %v", CC.SepOnly, tc.expSepOnly)
			}
			if len(CC.Source) != len(tc.expSource) {
				t.Fatalf("Source length = %d; want %d", len(CC.Source), len(tc.expSource))
			}
			for i := range CC.Source {
				if CC.Source[i] != tc.expSource[i] {
					t.Errorf("Source[%d] = %q; want %q", i, CC.Source[i], tc.expSource[i])
				}
			}
		})
	}
}

func TestValidateColumnsRange_Valid(t *testing.T) {
	valid := []string{
		"1,2,3", "1-5", "-3", "5-", "1,3-5,-7",
	}
	for _, f := range valid {
		CC.Fields = f
		if err := validateColumnsRange(); err != nil {
			t.Errorf("valid range %q returned error: %v", f, err)
		}
	}
}

func TestValidateColumnsRange_Invalid(t *testing.T) {
	invalid := []string{
		"a,b", "1,,2", "1-2-3", "1-2,k-",
	}
	for _, f := range invalid {
		CC.Fields = f
		if err := validateColumnsRange(); err == nil {
			t.Errorf("invalid range %q did not return error", f)
		}
	}
}

func TestPreprocessArgs_GroupedFlags(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	cases := []struct {
		name      string
		inputArgs []string
		expected  []string
	}{
		{
			name:      "grouped flags -sd",
			inputArgs: []string{"cmd", "-sd", "file.txt"},
			expected:  []string{"cmd", "-s", "-d", "file.txt"},
		},
		{
			name:      "single flags",
			inputArgs: []string{"cmd", "-s", "-d", "file.txt"},
			expected:  []string{"cmd", "-s", "-d", "file.txt"},
		},
		{
			name:      "mixed grouped and single",
			inputArgs: []string{"cmd", "-fsd", "input.txt"},
			expected:  []string{"cmd", "-f", "-s", "-d", "input.txt"},
		},
		{
			name:      "non-flag arguments untouched",
			inputArgs: []string{"cmd", "file1.txt", "file2.txt"},
			expected:  []string{"cmd", "file1.txt", "file2.txt"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			os.Args = tc.inputArgs
			preprocessArgs()

			if len(os.Args) != len(tc.expected) {
				t.Fatalf("len(os.Args) = %d; want %d", len(os.Args), len(tc.expected))
			}
			for i := range os.Args {
				if os.Args[i] != tc.expected[i] {
					t.Errorf("os.Args[%d] = %q; want %q", i, os.Args[i], tc.expected[i])
				}
			}
		})
	}
}
