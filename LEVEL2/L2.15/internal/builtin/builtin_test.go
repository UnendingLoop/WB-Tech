package builtin

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestEcho(t *testing.T) {
	var buf bytes.Buffer
	err := ECHO([]string{"echo", "hello", "world"}, nil, &buf)
	if err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(buf.String())
	if got != "hello world" {
		t.Errorf("expected 'hello world', got %q", got)
	}
}

func TestPWD(t *testing.T) {
	var buf bytes.Buffer
	err := PWD([]string{"pwd"}, nil, &buf)
	if err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(buf.String())
	wd, _ := os.Getwd()
	if got != wd {
		t.Errorf("expected %q, got %q", wd, got)
	}
}

func TestCD(t *testing.T) {
	orig, _ := os.Getwd()
	defer func() {
		if err := os.Chdir(orig); err != nil {
			t.Fatalf("Failed to change dir: %v", err)
		}
	}()
	tmp := os.TempDir()
	err := CD([]string{"cd", tmp}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	wd, _ := os.Getwd()
	if wd != tmp {
		t.Errorf("expected wd %q, got %q", tmp, wd)
	}
}

func TestIsBuiltin(t *testing.T) {
	if !IsBuiltin("echo") {
		t.Error("echo should be builtin")
	}
	if IsBuiltin("notbuiltin") {
		t.Error("notbuiltin should not be builtin")
	}
}

func TestPS(t *testing.T) {
	var buf bytes.Buffer
	err := PS([]string{"ps"}, nil, &buf)
	if err != nil {
		t.Fatalf("ps failed: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "PID") || !strings.Contains(out, "CMD") {
		t.Errorf("ps output missing headers: %q", out)
	}
	// Должен быть хотя бы один процесс (сам тест)
	if !strings.Contains(out, strconv.Itoa(os.Getpid())) {
		t.Errorf("ps output missing self pid: %q", out)
	}
}

func TestKill_InvalidArgs(t *testing.T) {
	var buf bytes.Buffer
	err := KILL([]string{"kill"}, nil, &buf)
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Errorf("expected usage error, got %v", err)
	}
}

func TestKill_InvalidPID(t *testing.T) {
	var buf bytes.Buffer
	err := KILL([]string{"kill", "notanumber"}, nil, &buf)
	if err == nil || !strings.Contains(err.Error(), "invalid pid") {
		t.Errorf("expected invalid pid error, got %v", err)
	}
}

func TestKill_NoSuchProcess(t *testing.T) {
	var buf bytes.Buffer
	// Используем заведомо несуществующий PID
	err := KILL([]string{"kill", "999999"}, nil, &buf)
	if err == nil || !strings.Contains(err.Error(), "process not found") {
		t.Errorf("expected process not found error, got %v", err)
	}
}
