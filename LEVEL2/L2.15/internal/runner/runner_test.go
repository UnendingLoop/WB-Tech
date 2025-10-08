//go:build !integration

package runner

import (
	"io"
	"miniShell/internal/builtin"
	"miniShell/internal/model"
	"os"
	"strings"
	"testing"
)

// mock builtin for testing
func mockEcho(args []string, in io.Reader, out io.Writer) error {
	for _, arg := range args[1:] {
		out.Write([]byte(arg + " "))
	}
	out.Write([]byte("\n"))
	return nil
}

func TestRunPipe_SingleBuiltin(t *testing.T) {
	builtin.BuiltInOps["echo"] = mockEcho
	p := model.Pipeline{
		Commands: []model.Command{
			{Args: []string{"echo", "hello", "world"}},
		},
	}
	// redirect stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := RunPipe(p)
	w.Close()
	os.Stdout = oldStdout
	if err != nil {
		t.Fatalf("RunPipe failed: %v", err)
	}
	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), "hello world") {
		t.Errorf("unexpected output: %s", string(out))
	}
}

func TestRunPipe_PipeBuiltinToBuiltin(t *testing.T) {
	builtin.BuiltInOps["echo"] = mockEcho
	builtin.BuiltInOps["upper"] = func(args []string, in io.Reader, out io.Writer) error {
		data, _ := io.ReadAll(in)
		out.Write([]byte(strings.ToUpper(string(data))))
		return nil
	}
	p := model.Pipeline{
		Commands: []model.Command{
			{Args: []string{"echo", "foo", "bar"}},
			{Args: []string{"upper"}},
		},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := RunPipe(p)
	w.Close()
	os.Stdout = oldStdout
	if err != nil {
		t.Fatalf("RunPipe failed: %v", err)
	}
	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), "FOO BAR") {
		t.Errorf("unexpected output: %s", string(out))
	}
}

func TestRunConditional_AndOr(t *testing.T) {
	builtin.BuiltInOps["true"] = func(args []string, in io.Reader, out io.Writer) error { return nil }
	builtin.BuiltInOps["false"] = func(args []string, in io.Reader, out io.Writer) error { return io.EOF }
	builtin.BuiltInOps["echo"] = mockEcho

	cond := model.Conditional{
		Pipelines: []model.Pipeline{
			{Commands: []model.Command{{Args: []string{"true"}}}},
			{Commands: []model.Command{{Args: []string{"echo", "ok"}}}},
			{Commands: []model.Command{{Args: []string{"false"}}}},
			{Commands: []model.Command{{Args: []string{"echo", "fail"}}}},
		},
		Operators: []string{"&&", "||", "&&"},
	}
	err := RunConditional(cond)
	if err != nil {
		t.Errorf("RunConditional failed: %v", err)
	}
}
