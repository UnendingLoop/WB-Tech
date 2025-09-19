// Package sorter contains the main 'management' logic for sorting input data
package sorter

import (
	"errors"
	"sort"
	"strconv"

	"sortClone/internal/model"
	"sortClone/internal/utils"
)

var DictMonths = map[string]string{"Jan": "01", "Feb": "02", "Mar": "03", "Apr": "04", "May": "05", "Jun": "06", "Jul": "07", "Aug": "08", "Sep": "09", "Oct": "10", "Nov": "11", "Dec": "12"}

// Sort - sorts input data according to flags fetched from cmd; returns sorted lines
func Sort(lines []string) []string {
	sort.Slice(lines, UniversalComparator(lines))
	// если есть флаг u, сразу убираем дубликаты; работает по колонке из k если он указан
	if model.OptsContainer.Unique {
		lines = utils.UniqueLines(lines)
	}
	return lines
}

// UniversalComparator - is a comparing function for sorting
func UniversalComparator(lines []string) func(i, j int) bool {
	return func(i, j int) bool {
		a, b := utils.GetKey(lines[i]), utils.GetKey(lines[j])
		// если есть флаг b, убираем хвостовые пробелы
		if model.OptsContainer.IgnSpaces {
			a, b = utils.RemoveTrailBlanks(a, b)
		}

		// парсим кило/мега/гига-байты
		if model.OptsContainer.HumanSort {
			ah, err0 := utils.ParseHumanSize(a)
			bh, err1 := utils.ParseHumanSize(b)

			humanResult := utils.SmallComparator(a, b, ah, bh, err0, err1)
			return ApplyReverse(humanResult, model.OptsContainer.Reverse)
		}

		// парсим месяцы; если месяц не определился - используем исходную строку и приоритет отдаем месяцу
		if model.OptsContainer.Monthly {
			am, oka := DictMonths[a[:3]]
			bm, okb := DictMonths[b[:3]]
			var errA, errB error
			if !oka {
				errA = errors.New("invalid month")
			}
			if !okb {
				errB = errors.New("invalid month")
			}
			monthlyResult := utils.SmallComparator(a, b, am, bm, errA, errB)
			return ApplyReverse(monthlyResult, model.OptsContainer.Reverse)
		}

		// парсим строки в число
		if model.OptsContainer.Numeric {
			ai, err2 := strconv.Atoi(a)
			bi, err3 := strconv.Atoi(b)

			numResult := utils.SmallComparator(a, b, ai, bi, err2, err3)
			return ApplyReverse(numResult, model.OptsContainer.Reverse)
		}

		if model.OptsContainer.Reverse {
			return b < a
		}
		return a < b
	}
}

func ApplyReverse(result, toReverse bool) bool {
	if toReverse {
		return !result
	}
	return result
}
