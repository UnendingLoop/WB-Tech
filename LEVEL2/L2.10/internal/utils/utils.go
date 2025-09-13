// Package utils provides small internal utilities
package utils

import (
	"strconv"
	"strings"

	"sortClone/internal/model"
)

var multipliers = map[string]float64{
	"K": 1024,
	"M": 1024 * 1024,
	"G": 1024 * 1024 * 1024,
	"T": 1024 * 1024 * 1024 * 1024,
}

// GetKey - returns value from the specified column of input line using specified delimeter
func GetKey(line string, column int, delim string) string {
	if column <= 0 {
		return line
	}

	parts := strings.Split(line, delim)
	if column-1 > len(parts) {
		return line
	}
	if column-1 < len(parts) {
		return parts[column-1]
	}

	return line
}

// UniqueLines - removes duplicate lines from input using a map
func UniqueLines(lines []string, opts model.Options) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, line := range lines {
		column := GetKey(line, opts.Column, opts.Delimeter)
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
