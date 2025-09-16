// Package reader provides means of reading input data
package reader

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/google/uuid"
)

const MaxFileSize int64 = 1024 * 1024 * 100

// ReadInput - reads input data that needs to be sorted whether it is a file or raw lines.
// Divides input file into tmp-smallfiles if it is too big
// Returns array of lines or array of filenames for further processing
func ReadInput(args []string) ([]string, []string, error) {
	var scanner *bufio.Scanner

	if len(args) > 1 {
		return nil, nil, errors.New("too many arguments to process")
	}
	// сначала смотрим свойства файла - открывается ли он и какой у него размер
	if len(args) == 1 {
		info, err := os.Stat(args[0])
		if err != nil || info.IsDir() { // если это папка или ошибка при открытии - возвращаем nil
			return nil, nil, errors.New("couldn't open specified input file")
		}
		// если размер файла более 100Мб - дробим на временные файлы
		if info.Size() > MaxFileSize {
			files, err := divideBigFile(args[0])
			return nil, files, err
		}

		file, err := os.Open(args[0])
		if err != nil {
			return nil, nil, err
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	if len(args) == 0 {
		scanner = bufio.NewScanner(os.Stdin)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	scanner.Buffer(make([]byte, 0, 1024*1024*5), 1024*1024*10)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil, scanner.Err()
}

func divideBigFile(filename string) ([]string, error) {
	resultList := []string{}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	line := ""
	chunk := strings.Builder{}

	for scanner.Scan() {
		line = scanner.Text() + "\n"
		chunk.WriteString(line)
		if chunk.Len() >= int(MaxFileSize) {
			tmpName, err := writeToTMP(chunk.String())
			if err != nil {
				return nil, err
			}
			resultList = append(resultList, tmpName)
			chunk.Reset()
		}
	}

	if chunk.Len() > 0 {
		tmpName, err := writeToTMP(chunk.String())
		if err != nil {
			return nil, err
		}
		resultList = append(resultList, tmpName)
	}
	return resultList, nil
}

func writeToTMP(input string) (string, error) {
	err := os.MkdirAll("./tmp", 0o755)
	if err != nil {
		return "", err
	}

	newName := "./tmp/" + uuid.New().String()
	file, err := os.Create(newName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.WriteString(file, input)
	if err != nil {
		return "", err
	}

	return newName, nil
}
