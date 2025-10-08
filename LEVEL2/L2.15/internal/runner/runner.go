// Package runner executes commands from model.Pipeline
package runner

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sync"

	"miniShell/internal/builtin"
	"miniShell/internal/model"
)

func RunPipe(p model.Pipeline) error {
	num := len(p.Commands)
	if num == 0 {
		return nil
	}

	readers := make([]io.Reader, num)
	writers := make([]io.Writer, num)
	procs := make([]*exec.Cmd, num)
	filesToClose := make([]*os.File, num)
	var wg sync.WaitGroup

	// Создаем пайпы для промежуточных команд
	for i := 0; i < num-1; i++ {
		r, w := io.Pipe()
		writers[i] = w
		readers[i+1] = r
	}
	writers[num-1] = os.Stdout
	readers[0] = os.Stdin

	// Открываем файлы для редиректов
	for i, cmd := range p.Commands {
		if cmd.InputFile != "" {
			f, err := os.Open(cmd.InputFile)
			if err != nil {
				return fmt.Errorf("cannot open input file %q: %w", cmd.InputFile, err)
			}
			readers[i] = f
			filesToClose[i] = f
		}
		if cmd.OutputFile != "" {
			flag := os.O_CREATE | os.O_WRONLY
			if cmd.Append {
				flag |= os.O_APPEND
			} else {
				flag |= os.O_TRUNC
			}
			f, err := os.OpenFile(cmd.OutputFile, flag, 0o644)
			if err != nil {
				return fmt.Errorf("cannot open output file %q: %w", cmd.OutputFile, err)
			}
			writers[i] = f
			filesToClose[i] = f
		}
	}

	// Запускаем команды
	for i, cmd := range p.Commands {
		stdin := readers[i]
		stdout := writers[i]
		cmd.Args = replaceEnvVars(cmd.Args)

		if builtinFn, ok := builtin.BuiltInOps[cmd.Args[0]]; ok {
			wg.Add(1)
			go func(fn func([]string, io.Reader, io.Writer) error, args []string, in io.Reader, out io.Writer, idx int) {
				defer wg.Done()
				if err := fn(args, in, out); err != nil {
					log.Printf("Error occured during %q: %v", cmd.Args[0], err)
				}
				// Закрываем PipeWriter сразу после завершения встроенной команды
				if idx < num-1 {
					if w, ok := out.(*io.PipeWriter); ok {
						w.Close()
					}
				}
			}(builtinFn, cmd.Args, stdin, stdout, i)
			continue
		}
		c := exec.Command(cmd.Args[0], cmd.Args[1:]...)
		c.Stdin = stdin
		c.Stdout = stdout
		c.Stderr = os.Stderr

		if err := c.Start(); err != nil {
			return fmt.Errorf("cannot start command %v: %w", cmd.Args, err)
		}
		procs[i] = c
	}

	wg.Wait()

	// Закрываем пайпы и файлы после Wait
	for i, c := range procs {
		if c == nil {
			continue // встроенная команда
		}
		if err := c.Wait(); err != nil {
			return fmt.Errorf("command %v failed: %w", c.Args, err)
		}
		// Закрываем PipeWriter, который был output для этой команды (если есть)
		if i < num-1 {
			if w, ok := writers[i].(*io.PipeWriter); ok {
				w.Close()
			}
		}
		// Закрываем input/output файл, если он был открыт
		if f := filesToClose[i]; f != nil && f != os.Stdout && f != os.Stdin && f != os.Stderr {
			f.Close()
		}
	}

	return nil
}

func RunConditional(cond model.Conditional) error {
	var lastErr error

	for i, p := range cond.Pipelines {
		err := RunPipe(p)
		success := err == nil

		if i < len(cond.Operators) {
			op := cond.Operators[i]
			if op == "&&" && !success {
				return err
			}
			if op == "||" && success {
				return nil
			}
		}
		lastErr = err
	}

	return lastErr
}

func replaceEnvVars(args []string) []string {
	res := make([]string, len(args))
	re := regexp.MustCompile(`\$(\w+)`) // ищем $VAR

	for i, arg := range args {
		res[i] = re.ReplaceAllStringFunc(arg, func(match string) string {
			varName := match[1:] // убираем $
			if val, ok := os.LookupEnv(varName); ok {
				return val
			}
			// если переменной нет — оставляем $VAR
			return match
		})
	}

	return res
}
