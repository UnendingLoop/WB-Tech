// Package printer provides method for printing input lines according to specified range of columns
package printer

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"cutClone/internal/parser"
)

func SplitAndPrintLine(line string) {
	buf := strings.Builder{}

	lineColumns := strings.Split(line, parser.CC.Delim)
	userColumns := strings.Split(parser.CC.Fields, ",")

	regex0 := regexp.MustCompile(`^\d+$`)        // одна колонка
	regex1 := regexp.MustCompile(`^\d+-$`)       // указанная колонка и до конца
	regex2 := regexp.MustCompile(`^-\d+$`)       // все колонки до указанного номера
	regex3 := regexp.MustCompile(`^\d+(-\d+)?$`) // диапазон колонок

	for _, v := range userColumns {
		switch {
		case regex0.MatchString(v): // одна колонка
			n, err := strconv.Atoi(v)
			if err != nil {
				log.Printf("Failed to convert 'N' column to int: %v", err)
			}
			buf.WriteString(formResultLine(n-1, n-1, lineColumns))

		case regex1.MatchString(v): // указанная колонка и до конца
			n, err := strconv.Atoi(strings.TrimRight(v, "-"))
			if err != nil {
				log.Printf("Failed to convert 'N-' column to int: %v", err)
			}
			buf.WriteString(formResultLine(n-1, len(lineColumns)-1, lineColumns))

		case regex2.MatchString(v): // все колонки до указанного номера
			n, err := strconv.Atoi(strings.TrimLeft(v, "-"))
			if err != nil {
				log.Printf("Failed to convert '-N' column to int: %v", err)
			}
			buf.WriteString(formResultLine(0, n-1, lineColumns))

		case regex3.MatchString(v): // диапазон колонок
			s1, s2, _ := strings.Cut(v, "-")
			m, err := strconv.Atoi(s1)
			if err != nil {
				log.Printf("Failed to convert 'M-N' M-column to int: %v", err)
			}
			n, err := strconv.Atoi(s2)
			if err != nil {
				log.Printf("Failed to convert 'M-N' N-column to int: %v", err)
			}
			buf.WriteString(formResultLine(m-1, n-1, lineColumns))
		}
	}

	if result := buf.String(); len(result) != 0 {
		fmt.Println(strings.TrimRight(result, parser.CC.Delim))
	}
}

func formResultLine(start, end int, lineArr []string) string {
	buf := strings.Builder{}

	if start == end {
		if start < len(lineArr) {
			return lineArr[start] + parser.CC.Delim
		}
		return ""
	}

	switch {
	case start > end:
		if end > len(lineArr)-1 {
			return ""
		}
		if start > len(lineArr)-1 {
			start = len(lineArr) - 1
		}

		for i := start; i >= end; i-- {
			buf.WriteString(lineArr[i] + parser.CC.Delim)
		}

	case start < end:
		if start > len(lineArr)-1 {
			return ""
		}
		if end > len(lineArr)-1 {
			end = len(lineArr) - 1
		}

		for i := start; i <= end; i++ {
			buf.WriteString(lineArr[i] + parser.CC.Delim)
		}
	}

	return buf.String()
}
