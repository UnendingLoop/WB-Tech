// Package reader provides means of reading input data
package reader

import (
	"bufio"
	"os"
)

// ReadInput - reads input data that needs to be sorted whether it is a file or raw lines.
func ReadInput(args []string) ([]string, error) {
	var scanner *bufio.Scanner
	if len(args) == 1 && fileExists(args[0]) {
		file, err := os.Open(args[0])
		if err != nil {
			return nil, err
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}
