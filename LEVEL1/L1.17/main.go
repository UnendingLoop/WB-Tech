package main

import (
	"fmt"
	"math/rand"
	"slices"
	"time"
)

/*
Реализовать алгоритм бинарного поиска встроенными методами языка.
Функция должна принимать отсортированный слайс и искомый элемент,
возвращать индекс элемента или -1, если элемент не найден.

Подсказка: можно реализовать рекурсивно или итеративно, используя цикл for.
*/
func generator(min, max, size int) []int {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	output := make([]int, size)
	for i := range size {
		output[i] = random.Intn(max-min+1) + min
	}

	return output
}

// Рекурсивный вариант
func binarySearchRecursive(arr []int, start, end, x int) int {
	if start > end {
		return -1
	}

	index := (start + end) / 2

	switch {
	case arr[index] == x:
		return index
	case arr[index] > x:
		return binarySearchRecursive(arr, start, index-1, x)
	default:
		return binarySearchRecursive(arr, index+1, end, x)
	}
}

// Итеративный вариант
func binarySearchIterative(arr []int, x int) int {
	start, end := 0, len(arr)-1

	for start <= end {
		mid := (start + end) / 2
		switch {
		case arr[mid] == x:
			return mid
		case arr[mid] > x:
			end = mid - 1
		default: //arr[mid] < x
			start = mid + 1
		}
	}
	return -1
}

func main() {
	//сгенерируем массив и отсортируем
	var a, b, c int
	fmt.Println("Specify positive values: \n-min, \n-max \n-size \nof array to generate using spaces as a separator:")
	fmt.Scanln(&a, &b, &c)
	if a >= b || c == 0 || a < 0 || b < 0 {
		fmt.Println("Invalid input data! Restart the app and try again.")
		return
	}
	inputArray1 := generator(a, b, c)
	slices.Sort(inputArray1)
	inputArray2 := make([]int, len(inputArray1))
	copy(inputArray2, inputArray1)
	fmt.Println("Generated sorted array is:\n", inputArray1)

	fmt.Println("What number would you like to find in the array?")
	x := 0
	fmt.Scanln(&x)

	index1 := binarySearchRecursive(inputArray1, 0, len(inputArray1)-1, x)
	index2 := binarySearchIterative(inputArray2, x)
	if index2 == index1 || index1 < 0 {
		fmt.Println("The specified number not found.")
		return
	}
	fmt.Println("Recursive binary search result:", index1)
	fmt.Println("Iterative binary search result:", index2)
}
