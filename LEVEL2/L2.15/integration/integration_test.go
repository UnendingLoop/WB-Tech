package integration

import (
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"miniShell/internal/parser"
	"miniShell/internal/runner"
)

func TestIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	tmp := t.TempDir()
	pwdFile := path.Join(tmp, "pwd.txt")

	os.Setenv("TEST_VAR", "HelloWorld")

	// Массив кейсов: команда, имя файла для проверки, ожидаемые строки
	type testCase struct {
		cmd      string
		outFile  string
		expected []string
	}

	outFile1 := path.Join(tmp, "out1.txt")
	outFile2 := path.Join(tmp, "out2.txt")
	outFile3 := path.Join(tmp, "out3.txt")
	tmpFile := path.Join(tmp, "tmp.txt")

	cases := []testCase{
		// echo и редирект
		{"echo First line > " + outFile1, outFile1, []string{"First line"}},
		{"echo $TEST_VAR >> " + outFile1, outFile1, []string{"HelloWorld"}},
		// пайп и echo
		{"echo hello world | tr a-z A-Z >> " + outFile2, outFile2, []string{"HELLO WORLD"}},
		// условный оператор &&
		{"echo line1 > " + tmpFile + " && echo line2 >> " + tmpFile, tmpFile, []string{"line1", "line2"}},
		// условный оператор || (не выполнится)
		{"false || echo should_run >> " + tmpFile, tmpFile, []string{"should_run"}},
		// cd и pwd
		{"cd " + tmpDir, "", nil},
		{"pwd > " + pwdFile, pwdFile, []string{tmpDir}},
	}

	for _, tc := range cases {
		cond := parser.ParseConditional(tc.cmd)
		if err := runner.RunConditional(cond); err != nil {
			t.Fatalf("Failed to run line %q: %v", tc.cmd, err)
		}
		time.Sleep(1 * time.Second) // чуть меньше, но всё равно даём системе время
		if tc.outFile != "" && len(tc.expected) > 0 {
			data, err := os.ReadFile(tc.outFile)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", tc.outFile, err)
			}
			str := strings.TrimSpace(string(data))
			for _, want := range tc.expected {
				if !strings.Contains(str, want) {
					t.Errorf("%s content incorrect, want %q in:\n%s", tc.outFile, want, str)
				}
			}
		}
	}

	// чистим временные файлы
	os.Remove(outFile1)
	os.Remove(outFile2)
	os.Remove(outFile3)
	os.Remove(tmpFile)
	os.Remove(pwdFile)
}
