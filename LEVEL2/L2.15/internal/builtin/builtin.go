// Package builtin provides built-in commands for minishell: cd, pwd, echo, kill, ps
package builtin

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

type BuiltinFunc func(args []string, stdin io.Reader, stdout io.Writer) error

var BuiltInOps = map[string]BuiltinFunc{
	"cd":   CD,
	"pwd":  PWD,
	"echo": ECHO,
	"ps":   PS,
	"kill": KILL,
}

// IsBuiltin проверяет, встроенная ли команда.
func IsBuiltin(name string) bool {
	_, ok := BuiltInOps[name]
	return ok
}

// RunBuiltin запускает встроенную команду.
func RunBuiltin(name string, args []string, stdin io.Reader, stdout io.Writer) error {
	if handler, ok := BuiltInOps[name]; ok {
		return handler(args, stdin, stdout)
	}
	return fmt.Errorf("unknown builtin command: %s", name)
}

func ECHO(args []string, stdin io.Reader, stdout io.Writer) error {
	if len(args) == 1 {
		_, err := fmt.Fprintln(stdout)
		return err
	}

	startIdx := 1
	newline := true
	if args[1] == "-n" {
		newline = false
		startIdx = 2
	}

	out := strings.Join(args[startIdx:], " ")
	if newline {
		_, err := fmt.Fprintln(stdout, out)
		return err
	} else {
		_, err := fmt.Fprint(stdout, out)
		return err
	}
}

func CD(args []string, stdin io.Reader, stdout io.Writer) error {
	switch len(args) {
	case 1:
		home := os.Getenv("HOME")
		if home == "" {
			return fmt.Errorf("HOME not set")
		}
		return os.Chdir(home)
	case 2:
		if args[1] == "-" {
			return os.Chdir("..")
		}
		return os.Chdir(args[1])
	default:
		return fmt.Errorf("usage: cd <path>")
	}
}

func PWD(args []string, stdin io.Reader, stdout io.Writer) error {
	if len(args) > 1 {
		return fmt.Errorf("usage: pwd")
	}
	curDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current dir: %v", err)
	}
	fmt.Fprintln(stdout, curDir)
	return nil
}

func PS(args []string, stdin io.Reader, stdout io.Writer) error {
	procs, err := process.Processes()
	if err != nil {
		return err
	}

	sort.Slice(procs, func(i, j int) bool {
		return procs[i].Pid < procs[j].Pid
	})

	fmt.Fprintf(stdout, "%-8s %-20s\n", "PID", "CMD")
	for _, p := range procs {
		name, err := p.Name()
		if err != nil {
			name = "?"
		}

		fmt.Fprintf(stdout, "%-8d %-20s\n", p.Pid, name)
	}
	return nil
}

func KILL(args []string, stdin io.Reader, stdout io.Writer) error {
	startIdx := 1
	sig := os.Kill // по умолчанию

	// На Unix можно парсить сигнал, на Windows — только Kill
	// Можно добавить build tag для разных платформ

	if len(args) <= startIdx {
		return fmt.Errorf("usage: kill <pid>")
	}

	pidStr := args[startIdx]
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return fmt.Errorf("invalid pid: %s", pidStr)
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("process not found: %v", err)
	}

	if err := proc.Signal(sig); err != nil {
		return fmt.Errorf("failed to kill process %d: %v", pid, err)
	}

	fmt.Fprintf(stdout, "Process %d terminated\n", pid)
	return nil
}
