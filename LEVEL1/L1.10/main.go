package main

import (
	"fmt"
	"slices"
	"strings"
)

/*
Дана последовательность температурных колебаний: -25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5.
Объединить эти значения в группы с шагом 10 градусов.
Пример: -20:{-25.4, -27.0, -21.0}, 10:{13.0, 19.0, 15.5}, 20:{24.5}, 30:{32.5}.
Пояснение: диапазон -20 включает значения от -20 до -29.9, диапазон 10 – от 10 до 19.9, и т.д.
Порядок в подмножествах не важен.
*/

func main() {
	input := []float32{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
	groups := make(map[int][]float32)
	for _, v := range input {
		gr := int(v) / 10
		groups[gr*10] = append(groups[gr*10], v)
	}

	//Приведение ответа к формату из задания
	groupsSlice := []string{}
	for k, v := range groups {
		str := fmt.Sprintf("%d:%.1f", k, v)
		str = strings.ReplaceAll(str, "[", "{")
		str = strings.ReplaceAll(str, "]", "}")
		groupsSlice = append(groupsSlice, str)
	}
	slices.Sort(groupsSlice)
	result := fmt.Sprint(groupsSlice)
	result = strings.ReplaceAll(result, "} ", "}, ")
	fmt.Println(result[1 : len(result)-1])
}
