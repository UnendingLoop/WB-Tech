// Package reader provides reading input lines from file(s) or stdIn, also implements check for flag -s
package reader

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"cutClone/internal/parser"
	"cutClone/internal/printer"
)

func ProcessInput(fileName string) error {
	var scanner *bufio.Scanner
	var file *os.File

	// инициализация источника строк
	switch fileName {
	case "":
		scanner = bufio.NewScanner(os.Stdin)
	default:
		// проверяем открывается ли файл, и не папка ли это
		info, err := os.Stat(fileName)
		if err != nil { // если ошибка при открытии - проброс ошибки выше
			return fmt.Errorf("error opening file %q: %q", fileName, err)
		}
		if info.IsDir() { // если это папка - проброс ошибки выше
			return fmt.Errorf("specified source filename %q is a directory", fileName)
		}

		// открываем файл для чтения
		file, err = os.Open(fileName)
		if err != nil {
			return fmt.Errorf("couldn't open file %q: %q", fileName, err)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	// читаем строки в цикле
	for scanner.Scan() {
		line := scanner.Text()
		switch parser.CC.SepOnly {
		case true:
			if strings.Contains(line, parser.CC.Delim) {
				printer.SplitAndPrintLine(line)
			}
		default:
			printer.SplitAndPrintLine(line)
		}
	}

	return scanner.Err()
}
