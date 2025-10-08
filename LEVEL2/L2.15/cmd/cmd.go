// Package cmd works as an entry point into app
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"

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

		err := runner.RunConditional(result)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}
}
