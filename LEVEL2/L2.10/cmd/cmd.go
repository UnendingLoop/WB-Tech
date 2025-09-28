// Package cmd reads os arguments and calls main sort function
package cmd

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	heaper "sortClone/internal/heap"
	"sortClone/internal/model"
	"sortClone/internal/reader"
	"sortClone/internal/sorter"
	"sortClone/internal/utils"

	"github.com/spf13/cobra"
)

var cobraFlagParser = &cobra.Command{
	Use:   "sortClone [flags] [filename/rawLines]",
	Short: "sortClone — аналог UNIX-утилиты sort",
	RunE:  ExecuteSort,
}

// чтобы тесты проводить вводим подменяемые переменные:
var (
	ReadInputFunc      = reader.ReadInput
	OutputDST          io.Writer
	IsSortedMessage    = "Входные данные уже отсортированы."
	IsNotSortedMessage = "Входные данные несортированы."
)

// определяем флаги и помещаем в контейнер
func init() {
	cobraFlagParser.Flags().IntVarP(&model.OptsContainer.Column, "key", "k", 0, "Сортировать по указанному номеру колонки/столбца в строке. По умолчанию '0' - работа с целой строкой")
	cobraFlagParser.Flags().BoolVarP(&model.OptsContainer.Numeric, "numeric", "n", false, "Сортировать по числам")
	cobraFlagParser.Flags().BoolVarP(&model.OptsContainer.Reverse, "reverse", "r", false, "Сортировать в обратном порядке")
	cobraFlagParser.Flags().BoolVarP(&model.OptsContainer.Unique, "unique", "u", false, "Выводить только уникальные строки")
	cobraFlagParser.Flags().StringVarP(&model.OptsContainer.Delimeter, "delimiter", "d", "\t", "Указать разделитель колонок")
	cobraFlagParser.Flags().BoolVarP(&model.OptsContainer.Monthly, "monthly", "M", false, "Сортировать по месяцам")
	cobraFlagParser.Flags().BoolVarP(&model.OptsContainer.IgnSpaces, "blanks", "b", false, "Игнорировать хвостовые пробелы")
	cobraFlagParser.Flags().BoolVarP(&model.OptsContainer.CheckIfSorted, "check", "c", false, "Проверка отсортированности")
	cobraFlagParser.Flags().BoolVarP(&model.OptsContainer.HumanSort, "human", "H", false, "Человекочитаемая сортировка - поддержка суффиксов T/Тб,Г/Гб,М/Мб,К/Кб в лат./кир.")
	cobraFlagParser.Flags().StringVarP(&model.OptsContainer.WriteToFile, "output to file", "o", "", "Запись результата сортировки в новый файл с указанным названием")
}

func ExecuteSort(cmd *cobra.Command, args []string) error {
	// определяем что подано на вход для обработки
	lines, filesArray, err := ReadInputFunc(args)

	// ветка обработки строк из памяти
	if lines != nil {
		return linesOutput(sortLines(lines))
	}

	// ветка работы с множеством tmp-файлов
	if filesArray != nil && !model.OptsContainer.CheckIfSorted {
		return processTmpFiles(filesArray)
	}
	// ветка работы с множеством tmp-файлов и флагом -с
	if filesArray != nil && model.OptsContainer.CheckIfSorted {
		result, err := checkTmpFiles(filesArray)
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	}

	return err
}

func sortLines(lines []string) []string {
	if len(lines) <= 1 {
		return lines
	}

	// если есть флаг -c, проверка если строки уже отсортированы; работает по колонке из k если он указан
	if model.OptsContainer.CheckIfSorted {
		result := sort.SliceIsSorted(lines, sorter.UniversalComparator(lines))
		if result {
			fmt.Println(IsSortedMessage)
			return nil
		}
		fmt.Println(IsNotSortedMessage)
		return nil
	}

	// сама сортировка строк согласно всем остальным флагам
	return sorter.Sort(lines)
}

func linesOutput(lines []string) error {
	if lines == nil {
		return nil
	}
	// готовим выходной поток
	if model.OptsContainer.WriteToFile != "" {
		if OutputDST == nil {
			if model.OptsContainer.WriteToFile != "" {
				f, err := os.Create(model.OptsContainer.WriteToFile)
				if err != nil {
					return err
				}
				defer f.Close()
				OutputDST = f
			} else {
				OutputDST = os.Stdout
			}
		}
	} else {
		if OutputDST == nil {
			OutputDST = os.Stdout
		}
	}
	// обработка результата - или вывод на экран, или запись в файл
	for _, line := range lines {
		fmt.Fprintln(OutputDST, line)
	}

	return nil
}

func processTmpFiles(tmpfiles []string) error {
	// сортируем каждый временный файл
	for _, tmpName := range tmpfiles {
		lines, err := utils.ReadFileToRAM(tmpName)
		if err != nil {
			return err
		}
		result := sortLines(lines)
		err = utils.WriteLinesToFile(result, tmpName)
		if err != nil {
			return err
		}
	}
	// вызов мерджера временных файлов
	return mergeTmpFiles(tmpfiles)
}

func checkTmpFiles(tmpfiles []string) (string, error) {
	// проверяем все временные файлы в потоке по очереди
	lastline := ""
	started := false
	for _, tmpName := range tmpfiles {
		tmpLines, err := utils.ReadFileToRAM(tmpName)
		if err != nil {
			return "", err
		}
		for _, line := range tmpLines {
			if !started {
				lastline = line
				started = true
				continue
			}
			result := sort.SliceIsSorted([]string{lastline, line}, sorter.UniversalComparator([]string{lastline, line}))
			if !result {
				return IsNotSortedMessage, nil
			}
			lastline = line
		}
	}
	return IsSortedMessage, nil
}

func mergeTmpFiles(tmpfiles []string) error {
	tmpHeap := heaper.StrHeap{}
	heap.Init(&tmpHeap)
	// открываем временные файлы и добавляем по 1 элементу из каждого из них
	scanners := make([]*bufio.Scanner, len(tmpfiles))
	tmpFiles := make([]*os.File, len(tmpfiles))
	for i, fileName := range tmpfiles {
		tmpFile, err := os.Open(fileName)
		if err != nil {
			return fmt.Errorf("failed to open tmp-file %q for merging: %v", fileName, err)
		}
		tmpFiles[i] = tmpFile
		scanner := bufio.NewScanner(tmpFile)
		scanner.Buffer(make([]byte, 0, 1024*1024*5), 1024*1024*10)
		scanners[i] = scanner
		if scanners[i].Scan() {
			heap.Push(&tmpHeap, &heaper.FileLine{Value: scanner.Text(), FileID: i})
		}
	}
	defer func() {
		// закрываем все временные файлы
		for _, f := range tmpFiles {
			if err := f.Close(); err != nil {
				log.Fatalf("Failed to close tmp-file %q after merging: %v", f.Name(), err)
			}
		}
		// удаляем их
		if err := os.RemoveAll("tmp"); err != nil {
			log.Fatalf("Failed to remove tmp-directory: %v", err)
		}
	}()

	// готовим выходной поток
	if model.OptsContainer.WriteToFile != "" {
		if OutputDST == nil {
			if model.OptsContainer.WriteToFile != "" {
				f, err := os.Create(model.OptsContainer.WriteToFile)
				if err != nil {
					return err
				}
				defer f.Close()
				OutputDST = f
			} else {
				OutputDST = os.Stdout
			}
		}
	} else {
		if OutputDST == nil {
			OutputDST = os.Stdout
		}
	}
	// Достаем из кучи элементы
	lastWritten := ""
	for tmpHeap.Len() > 0 {
		item := heap.Pop(&tmpHeap).(*heaper.FileLine)
		switch model.OptsContainer.Unique { // ветвление в зависимости от флага Unique
		case true:
			if item.Value != lastWritten {
				if _, err := io.WriteString(OutputDST, item.Value+"\n"); err != nil {
					return err
				}
			}
		case false:
			if _, err := io.WriteString(OutputDST, item.Value+"\n"); err != nil {
				return err
			}
		}

		if scanners[item.FileID].Scan() {
			heap.Push(&tmpHeap, &heaper.FileLine{Value: scanners[item.FileID].Text(), FileID: item.FileID})
		}
		lastWritten = item.Value
	}
	// проверяем ошибки в сканерах
	for _, scanner := range scanners {
		if err := scanner.Err(); err != nil {
			return err
		}
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

// Execute - reads flags and launches sort function with flags
func Execute() {
	preprocessArgs()
	if err := cobraFlagParser.Execute(); err != nil {
		os.Exit(1)
	}
}
