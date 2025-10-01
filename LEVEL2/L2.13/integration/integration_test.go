package integration

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var testFile1Path, testFile2Path string

func TestMain(m *testing.M) {
	content1 := "0101\t0202\t0303\t0404\nAAAA\tBBBB\tCCCC\tDDDD\n1010\t1111\t1212\t1313"
	content2 := "0101\t0202\t0303\t0404\nEEEE\tFFFF\tGGGG\tHHHH\n1010\t1111\t1212\t1313"

	var err error
	testFile1Path, err = createTempFile(content1)
	if err != nil {
		panic(err)
	}
	testFile2Path, err = createTempFile(content2)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	os.Remove(testFile1Path)
	os.Remove(testFile2Path)
	os.Exit(code)
}

func createTempFile(content string) (string, error) {
	f, err := os.CreateTemp("", "cutClone_test_*.txt")
	if err != nil {
		return "", err
	}
	if _, err := f.WriteString(content); err != nil {
		return "", err
	}
	f.Close()
	return f.Name(), nil
}

func runCutClone(args ...string) (string, error) {
	cmd := exec.Command("go", append([]string{"run", "../main.go"}, args...)...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func TestIntegration_AllCases(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "basic_range_2-3",
			args: []string{"-f", "2-3", testFile1Path},
			want: []string{"0202\t0303", "BBBB\tCCCC", "1111\t1212"},
		},
		{
			name: "range_-3_(from_start)",
			args: []string{"-f", "-3", testFile1Path},
			want: []string{"0101\t0202\t0303", "AAAA\tBBBB\tCCCC", "1010\t1111\t1212"},
		},
		{
			name: "range_2-_(to_end)",
			args: []string{"-f", "2-", testFile1Path},
			want: []string{"0202\t0303\t0404", "BBBB\tCCCC\tDDDD", "1111\t1212\t1313"},
		},
		{
			name: "-s_flag_filters_lines_without_delimiter",
			args: []string{"-f", "1", "-s", testFile1Path},
			want: []string{"0101", "AAAA", "1010"}, // строки без '\t' не выводятся
		},
		{
			name: "multiple_files",
			args: []string{"-f", "1", testFile1Path, testFile2Path},
			want: []string{"0101", "AAAA", "1010", "0101", "EEEE", "1010"},
		},
		{
			name: "repeated_columns",
			args: []string{"-f", "1,1,2", testFile1Path},
			want: []string{"0101\t0101\t0202", "AAAA\tAAAA\tBBBB", "1010\t1010\t1111"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := runCutClone(tt.args...)
			if err != nil {
				t.Fatalf("command failed: %v\noutput:\n%s", err, out)
			}

			for _, wantLine := range tt.want {
				if !strings.Contains(out, wantLine) {
					t.Errorf("output missing expected line %q\nfull output:\n%s", wantLine, out)
				}
			}
		})
	}
}
