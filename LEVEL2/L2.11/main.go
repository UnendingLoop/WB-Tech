package main

import (
	"fmt"
	"slices"
	"strings"
)

func findAnagrams(input []string) map[string][]string {
	sortedAnaMap := make(map[string][]string)
	for _, v := range input {
		v = strings.ToLower(v)
		runes := []rune(v)
		slices.Sort(runes)
		sortedAnaMap[string(runes)] = append(sortedAnaMap[string(runes)], v)
	}

	resultAnaMap := make(map[string][]string)
	for _, v := range sortedAnaMap {
		if len(v) == 1 {
			continue
		}
		resultAnaMap[v[0]] = v
	}

	return resultAnaMap
}

func main() {
	inputStrArr := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	result := findAnagrams(inputStrArr)
	fmt.Println(result)
}
