package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
Реализовать пересечение двух неупорядоченных множеств (например, двух слайсов) — т.е.
вывести элементы, присутствующие и в первом, и во втором.

Пример:
A = {1,2,3}
B = {2,3,4}
Пересечение = {2,3}
*/

func main() {
	var a, b string
	var err error

	r := bufio.NewReader(os.Stdin)
	fmt.Println("Input the 1st array of int-numbers using 'space' as a separator in one line:")
	if a, err = r.ReadString('\n'); err != nil {
		fmt.Printf("Invalid input: %v. Try again.", err)
		return
	}
	fmt.Println("Input the 2nd array of int-numbers using 'space' as a separator in one line:")
	if b, err = r.ReadString('\n'); err != nil {
		fmt.Printf("Invalid input: %v. Try again.", err)
		return
	}

	aInts, err := sliceStrToSliceInt(strings.Fields(a))
	if err != nil {
		fmt.Printf("Failed to process input data: %v. Try again.", err)
		return
	}
	bInts, err := sliceStrToSliceInt(strings.Fields(b))
	if err != nil {
		fmt.Printf("Failed to process input data: %v. Try again.", err)
		return
	}

	result := findCommon(aInts, bInts)

	fmt.Println("Common elements are: ", result)

}

func sliceStrToSliceInt(input []string) ([]int, error) {
	result := make([]int, 0, len(input))
	for _, v := range input {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		result = append(result, n)
	}
	return result, nil
}

func findCommon(a, b []int) []int {
	if len(a) > len(b) {
		a, b = b, a // меняем местами, чтобы a всегда было меньше - экономия размера карты
	}

	result := []int{}
	aMap := make(map[int]struct{})
	for _, k := range a {
		aMap[k] = struct{}{}
	}

	rMap := make(map[int]struct{})
	for _, k := range b {
		if _, ok := aMap[k]; ok {
			if _, ok := rMap[k]; !ok { //избавляемся от дубликатов
				result = append(result, k)
				rMap[k] = struct{}{}
			}
		}
	}

	return result
}
