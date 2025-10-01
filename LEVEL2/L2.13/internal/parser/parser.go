// Package parser puts os.Args into SearchParam structure and validates it for any issues
package parser

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"cutClone/internal/model"
)

var CC model.CutConfig

func InitSearchParam() error {
	flagParser := flag.NewFlagSet("cutClone", flag.ExitOnError)

	f := flagParser.String("f", "", "specify columns to print like '1,3-5'")
	d := flagParser.String("d", "\t", "set user-specified delimeter('\t' by default)")
	s := flagParser.Bool("s", false, "show lines only containing user-specified delimeter(uses '\t' if -d is not set)")

	preprocessArgs()

	if err := flagParser.Parse(os.Args[1:]); err != nil {
		return err
	}

	CC = model.CutConfig{
		Fields:  *f,
		Delim:   *d,
		SepOnly: *s,
	}

	// Разбираемся с входом
	switch len(flagParser.Args()) {
	case 0:
		return errors.New("input not specified!\nUsage: grepClone [flags] pattern [file...]")
	default:
		CC.Source = flagParser.Args()
	}
	// проверка корректности указанного диапазона колонок и разворот
	if err := validateColumnsRange(); err != nil {
		return err
	}

	return nil
}

func preprocessArgs() {
	var args []string
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && len(arg) > 2 {
			// разбиваем группированные флаги если есть
			for _, ch := range arg[1:] {
				args = append(args, "-"+string(ch))
			}
		} else {
			args = append(args, arg)
		}
	}
	os.Args = args
}

func validateColumnsRange() error {
	re := regexp.MustCompile(`^(\d+(-\d+)?|\d+-|-\d+)(,(\d+(-\d+)?|\d+-|-\d+))*$`)
	if !re.MatchString(CC.Fields) {
		return fmt.Errorf("specified column range is incorrect")
	}
	return nil
}
