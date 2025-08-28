package main

import "fmt"

/*
Удалить i-ый элемент из слайса.
Продемонстрируйте корректное удаление без утечки памяти.

Подсказка: можно
- сдвинуть хвост слайса на место удаляемого элемента (copy(slice[i:], slice[i+1:])) и
- уменьшить длину слайса на 1.
*/

// removeWithKeepingOrder - исчключение iго элемента с возвратом нового слайса
func removeWithKeepingOrder[T any](arr []T, i int) []T {
	arr = append(arr[:i], arr[i+1:]...)
	b := make([]T, len(arr))
	n := copy(b, arr)
	//обнуляем пару элементов, чтобы проверить, что новый слайс уже не связан с исходным
	arr[0] = *new(T)
	arr[3] = *new(T)
	fmt.Println("How many elements copied:", n)
	fmt.Println("Array A:", arr)
	fmt.Println("Array B:", b)
	return b
}

// removeWithoutOrder - копирует значение последнего элемента в iй и создает копию этого слайса. Исходный слайс будет удален GC если на него более нет ссылок.
func removeWithoutOrder[T any](arr []T, i int) []T {
	arr[i] = arr[len(arr)-1] //свопаем iй с последним
	b := make([]T, len(arr)-1)
	n := copy(b, arr)
	//обнуляем пару элементов, чтобы проверить, что новый слайс уже не связан с исходным
	arr[0] = *new(T)
	arr[3] = *new(T)
	fmt.Println("How many elements copied:", n)
	fmt.Println("Array A:", arr)
	fmt.Println("Array B:", b)
	return b
}

// removeKeepingOrderWithoutCreation - исключает iй элемент из исходного массива посредством сдвига хвоста без создания копии; последний элемент обнуляется согласно типу, чтобы GC мог его очистить если это ссылочный элемент
func removeKeepingOrderWithoutCreation[T any](arr []T, i int) []T {
	copy(arr[i:], arr[i+1:])  // сдвигаем хвост влево
	arr[len(arr)-1] = *new(T) // обнуляем последний элемент
	return arr[:len(arr)-1]   // уменьшаем длину
}

func main() {
	input1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	result1 := removeWithKeepingOrder(input1, 3)

	input2 := []int{5, 6, 7, 8, 9, 10, 11}
	result2 := removeWithoutOrder(input2, 4)

	input3 := []int{10, 11, 12, 13, 14, 15, 16}
	result3 := removeKeepingOrderWithoutCreation(input3, 4)

	fmt.Println("\nRemoving with keeping order and creating new slice:", result1)
	fmt.Println("Without order with creating new slice:", result2)
	fmt.Println("Without creating new slice:", result3)

}
