// Package cmd reads os arguments and calls main sort function
package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"sortClone/internal/model"
	"sortClone/internal/reader"
	"sortClone/internal/sorter"

	"github.com/spf13/cobra"
)

var (
	column          int
	numeric         bool
	reverse         bool
	unique          bool
	delimeter       string
	monthly         bool
	ignoreEndSpaces bool
	checkSorted     bool
	humanSort       bool
	writeToFile     string
)

var rootCmd = &cobra.Command{
	Use:   "sortClone [flags] [filename/rawLines]",
	Short: "sortClone — аналог UNIX-утилиты sort",
	RunE:  sortData,
}

func init() {
	rootCmd.Flags().IntVarP(&column, "key", "k", 0, "Сортировать по указанному номеру колонки/столбца в строке. По умолчанию '0' - работа с целой строкой")
	rootCmd.Flags().BoolVarP(&numeric, "numeric", "n", false, "Сортировать по числам")
	rootCmd.Flags().BoolVarP(&reverse, "reverse", "r", false, "Сортировать в обратном порядке")
	rootCmd.Flags().BoolVarP(&unique, "unique", "u", false, "Выводить только уникальные строки")
	rootCmd.Flags().StringVarP(&delimeter, "delimeter", "d", "\t", "Указать разделитель колонок")
	rootCmd.Flags().BoolVarP(&monthly, "monthly", "M", false, "Сортировать по месяцам")
	rootCmd.Flags().BoolVarP(&ignoreEndSpaces, "blanks", "b", false, "Игнорировать хвостовые пробелы")
	rootCmd.Flags().BoolVarP(&checkSorted, "check", "c", false, "Проверка отсортированности")
	rootCmd.Flags().BoolVarP(&humanSort, "human", "H", false, "Человекочитаемая сортировка - поддержка суффиксов Тб,Гб,Мб,Кб")
	rootCmd.Flags().StringVarP(&writeToFile, "write2file", "w", "", "Запись результата сортировки в новый файл с указанным названием")
}

func sortData(cmd *cobra.Command, args []string) error {
	lines, err := reader.ReadInput(args)
	if err != nil {
		return err
	}

	parsedOptions := model.Options{
		Column:        column,
		Numeric:       numeric,
		Reverse:       reverse,
		Unique:        unique,
		Delimeter:     delimeter,
		Monthly:       monthly,
		IgnSpaces:     ignoreEndSpaces,
		CheckIfSorted: checkSorted,
		HumanSort:     humanSort,
		WriteToFile:   writeToFile,
	}

	// если есть флаг -c, проверка если строки уже отсортированы; работает по колонке из k если он указан
	if parsedOptions.CheckIfSorted {
		result := sort.SliceIsSorted(lines, sorter.UniversalComparator(parsedOptions, lines))
		if result {
			fmt.Println("Входные данные уже отсортированы.")
			return nil
		}
		fmt.Println("Входные данные несортированы.")
		return nil
	}
	// сама сортировка согласно всем остальным флагам
	result := sorter.Sort(lines, parsedOptions)

	// обработка результата - или вывод на экран, или запись в файл
	if writeToFile != "" {
		f, err := os.Create(writeToFile)
		if err != nil {
			return fmt.Errorf("не удалось создать файл: %w", err)
		}
		defer f.Close()

		for _, line := range result {
			fmt.Fprintln(f, line)
		}
	} else {
		for _, line := range result {
			fmt.Println(line)
		}
	}

	return nil
}

func preprocessArgs() {
	var args []string
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && len(arg) > 2 {
			// Разбиваем группированые флаги
			for _, ch := range arg[1:] {
				args = append(args, "-"+string(ch))
			}
		} else {
			args = append(args, arg)
		}
	}
	os.Args = args
}

// Execute - reads flags and launches sort function with flags
func Execute() {
	preprocessArgs()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
