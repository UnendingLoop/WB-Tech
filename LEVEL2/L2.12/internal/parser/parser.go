// Package parser puts os.Args into SearchParam structure and validates it for any issues
package parser

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"grepClone/internal/model"
)

var SP model.SearchParam

func InitSearchParam() error {
	preprocessArgs()

	flagParser := flag.NewFlagSet("grepClone", flag.ExitOnError)

	a := flagParser.Int("A", -1, "show N lines after target line")
	b := flagParser.Int("B", -1, "show N lines before target line")
	c := flagParser.Int("C", 0, "show N lines before and after target line(same as '-A N' and '-B N')")
	d := flagParser.Bool("c", false, "show only total number of matching lines(A/B/C/n are ignored)")
	e := flagParser.Bool("i", false, "all input lines will be lower-cased for search as well as the pattern itself")
	f := flagParser.Bool("v", false, "search only lines that DON'T match the specified pattern")
	g := flagParser.Bool("F", false, "specified pattern will be used strictly as a string, not regexp")
	h := flagParser.Bool("n", false, "enumerates output lines according to their order in input")

	if err := flagParser.Parse(os.Args[1:]); err != nil {
		return err
	}

	SP = model.SearchParam{
		CtxAfter:     *a,
		CtxBefore:    *b,
		CtxCircle:    *c,
		CountFound:   *d,
		IgnoreCase:   *e,
		InvertResult: *f,
		ExactMatch:   *g,
		EnumLine:     *h,
	}

	// Выравниваем значения контекста A и B по значению C
	if SP.CtxCircle > 0 {
		if SP.CtxAfter == -1 {
			SP.CtxAfter = SP.CtxCircle
		}
		if SP.CtxBefore == -1 {
			SP.CtxBefore = SP.CtxCircle
		}
	}
	if SP.CtxAfter == -1 {
		SP.CtxAfter = 0
	}
	if SP.CtxBefore == -1 {
		SP.CtxBefore = 0
	}

	// Приводим паттерн к нижнему регистру если стоят флаги 'F' и 'i'
	if SP.IgnoreCase && SP.ExactMatch {
		SP.Pattern = strings.ToLower(SP.Pattern)
	}

	// Разбираемся с паттерном и входом
	switch len(flagParser.Args()) {
	case 0:
		return errors.New("pattern not specified!\nUsage: grepClone [flags] pattern [file...]")
	case 1:
		SP.Pattern = flagParser.Args()[0]
	default:
		SP.Pattern = flagParser.Args()[0]
		SP.Source = flagParser.Args()[1:]
	}

	// ставим флаг чтобы печатать имя файла перед каждой строкой/суммой строк, если файлов несколько
	if len(SP.Source) > 1 {
		SP.PrintFileName = true
	}

	// сразу првоеряем корректность регулярки, если флаг F неактивен
	if _, err := regexp.Compile(SP.Pattern); !SP.ExactMatch && err != nil {
		return fmt.Errorf("specified pattern %q is incorrect regexp", SP.Pattern)
	}

	return nil
}

func preprocessArgs() {
	var args []string
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && len(arg) > 2 {
			// разбиваем группированые флаги если есть
			for _, ch := range arg[1:] {
				args = append(args, "-"+string(ch))
			}
		} else {
			args = append(args, arg)
		}
	}
	os.Args = args
}
