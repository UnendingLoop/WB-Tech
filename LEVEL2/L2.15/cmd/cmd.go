// Package cmd works as an entry point into app
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"

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
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "< ") {
			scriptFile := strings.TrimSpace(line[2:])
			if err := runScriptFile(scriptFile); err != nil {
				fmt.Fprintf(os.Stderr, "Script error: %v\n", err)
			}
			continue
		}
		result := parser.ParseConditional(line)

		err := runner.RunConditional(result)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}
}

func runScriptFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("cannot open script file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		cond := parser.ParseConditional(line)
		if err := runner.RunConditional(cond); err != nil {
			fmt.Fprintf(os.Stderr, "Error in script: %v\n", err)
		}
	}
	return scanner.Err()
}
