package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
Реализовать алгоритм быстрой сортировки массива встроенными средствами языка. Можно использовать рекурсию.

Подсказка: напишите функцию quickSort([]int) []int которая сортирует срез целых чисел.
Для выбора опорного элемента можно взять середину или первый элемент.
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

func qsort(arr []int, start, end int) {
	if start < end {
		pivot := partitioner(arr, start, end)
		qsort(arr, start, pivot-1)
		qsort(arr, pivot+1, end)
	}
}
func partitioner(arr []int, start, end int) int {
	pivot := arr[end] //берем последний эелемент в качестве опорного
	i := start - 1

	for j := start; j < end; j++ {
		if arr[j] < pivot {
			i++
			arr[j], arr[i] = arr[i], arr[j]
		}
	}
	arr[i+1], arr[end] = arr[end], arr[i+1]

	return i + 1
}

func main() {
	//сгенерируем массив
	var a, b, c int
	fmt.Println("Specify values: \n-min, \n-max \n-size \nof array to generate using spaces as a separator:")
	fmt.Scanln(&a, &b, &c)
	if a >= b || c == 0 {
		fmt.Println("Invalid input data! Restart the app and try again.")
		return
	}
	inputArray := generator(a, b, c)
	fmt.Println("Generated unsorted array is:\n", inputArray)

	qsort(inputArray, 0, len(inputArray)-1)

	fmt.Println("Qsort result is:\n", inputArray)
}
