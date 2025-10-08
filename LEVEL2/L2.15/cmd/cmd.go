// Package cmd works as an entry point into app
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"regexp"

	"miniShell/internal/parser"
	"miniShell/internal/runner"
)

func InitReadAndRun() {
	// Канал для Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	// Горутина для перехвата Ctrl+C
	go func() {
		for range sigCh {
			fmt.Print("\n> ") // просто выводим новый prompt
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() { // ловим Ctrl+D (Сtrl+Z на винде)
			fmt.Println("\nExiting shell…")
			break
		}
		line := scanner.Text()
		if line == "" {
			continue
		}
		result := parser.ParseConditional(line)

		// Подставляем переменные окружения в каждой команде каждого пайпа
		for i := range result.Pipelines {
			for j := range result.Pipelines[i].Commands {
				cmd := &result.Pipelines[i].Commands[j]
				cmd.Args = replaceEnvVars(cmd.Args)
			}
		}

		err := runner.RunConditional(result)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}
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
