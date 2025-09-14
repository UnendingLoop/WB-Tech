// Package utils provides small internal utilities
package utils

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"sortClone/internal/model"
)

var multipliers = map[string]float64{
	"K": 1024,
	"M": 1024 * 1024,
	"G": 1024 * 1024 * 1024,
	"T": 1024 * 1024 * 1024 * 1024,
	"К": 1024,
	"М": 1024 * 1024,
	"Г": 1024 * 1024 * 1024,
	"Т": 1024 * 1024 * 1024 * 1024,
}

// GetKey - returns value from the specified column of input line using specified delimeter
func GetKey(line string) string {
	if model.OptsContainer.Column <= 0 {
		return line
	}

	parts := strings.Split(line, model.OptsContainer.Delimeter)
	if model.OptsContainer.Column-1 > len(parts) {
		return line
	}
	if model.OptsContainer.Column-1 < len(parts) {
		return parts[model.OptsContainer.Column-1]
	}

	return line
}

// UniqueLines - removes duplicate lines from input using a map
func UniqueLines(lines []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, line := range lines {
		column := GetKey(line)
		if _, ok := seen[column]; !ok {
			seen[column] = struct{}{}
			result = append(result, line)
		}
	}
	return result
}

// RemoveTrailBlanks - removes all spaces from the end of every input line
func RemoveTrailBlanks(a, b string) (string, string) {
	return strings.TrimRight(a, " "), strings.TrimRight(b, " ")
}

// ParseHumanSize - parses values like 1024K/64M to make them comparable
func ParseHumanSize(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}

	// отделяем число от суффикса
	numPart := ""
	unitPart := ""
	for i, r := range s {
		if (r < '0' || r > '9') && r != '.' {
			numPart = s[:i]
			unitPart = strings.ToUpper(strings.TrimSpace(s[i:]))
			break
		}
	}
	if numPart == "" {
		numPart = s
	}

	value, err := strconv.ParseFloat(numPart, 64)
	if err != nil {
		return 0, err
	}

	if mul, ok := multipliers[unitPart]; ok {
		value *= mul
	}

	return value, nil
}

// SmallComparator - used for comparing 2 input lines according with result of their conversion to destination type
// Priority is given to successfully converted value, otherwise it will be compared as a simple string
func SmallComparator[T int | float64 | string](strA, strB string, a, b T, err1, err2 error) bool {
	switch {
	case err1 != nil && err2 != nil:
		return strA < strB
	case err1 != nil:
		return false
	case err2 != nil:
		return true
	default:
		return a < b
	}
}

func ReadFileToRAM(fileName string) ([]string, error) {
	var scanner *bufio.Scanner

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error while closing file '%v': %v", fileName, err)
		}
	}()

	scanner = bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	scanner.Buffer(make([]byte, 0, 1024*1024*5), 1024*1024*10)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func WriteLinesToFile(lines []string, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error while closing file '%v': %v", fileName, err)
		}
	}()

	builder := strings.Builder{}
	for _, line := range lines {
		builder.WriteString(line + "\n")
	}
	_, err = io.WriteString(file, builder.String())
	if err != nil {
		return err
	}
	return nil
}
