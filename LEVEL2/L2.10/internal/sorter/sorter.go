// Package sorter contains the main 'management' logic for sorting input data
package sorter

import (
	"errors"
	"sort"
	"strconv"

	"sortClone/internal/model"
	"sortClone/internal/utils"
)

var dictMonths = map[string]string{"Jan": "01", "Feb": "02", "Mar": "03", "Apr": "04", "May": "05", "Jun": "06", "Jul": "07", "Aug": "08", "Sep": "09", "Oct": "10", "Nov": "11", "Dec": "12"}

// Sort - sorts input data according to flags fetched from cmd; returns sorted lines
func Sort(lines []string, opts model.Options) []string {
	sort.Slice(lines, UniversalComparator(opts, lines))
	// если есть флаг u, сразу убираем дубликаты; работает по колонке из k если он указан
	if opts.Unique {
		lines = utils.UniqueLines(lines, opts)
	}
	return lines
}

// UniversalComparator - is a comparing function for sorting
func UniversalComparator(opts model.Options, lines []string) func(i, j int) bool {
	return func(i, j int) bool {
		a, b := utils.GetKey(lines[i], opts.Column, opts.Delimeter), utils.GetKey(lines[j], opts.Column, opts.Delimeter)

		// если есть флаг b, убираем хвостовые пробелы
		if opts.IgnSpaces {
			a, b = utils.RemoveTrailBlanks(a, b)
		}

		// парсим кило/мега/гига-байты
		if opts.HumanSort {
			ah, err0 := utils.ParseHumanSize(a)
			bh, err1 := utils.ParseHumanSize(b)

			humanResult := utils.SmallComparator(a, b, ah, bh, err0, err1)
			return applyReverse(humanResult, opts.Reverse)
		}

		// парсим месяцы; если месяц не определился - используем исходную строку и приоритет отдаем месяцу
		if opts.Monthly {
			am, oka := dictMonths[a]
			bm, okb := dictMonths[b]
			var errA, errB error
			if !oka {
				errA = errors.New("invalid month")
			}
			if !okb {
				errB = errors.New("invalid month")
			}
			monthlyResult := utils.SmallComparator(a, b, am, bm, errA, errB)
			return applyReverse(monthlyResult, opts.Reverse)
		}

		// парсим строки в число
		if opts.Numeric {
			ai, err2 := strconv.Atoi(a)
			bi, err3 := strconv.Atoi(b)

			numResult := utils.SmallComparator(a, b, ai, bi, err2, err3)
			return applyReverse(numResult, opts.Reverse)
		}

		if opts.Reverse {
			return b < a
		}
		return a < b
	}
}

func applyReverse(result, toReverse bool) bool {
	if toReverse {
		return !result
	}
	return result
}
