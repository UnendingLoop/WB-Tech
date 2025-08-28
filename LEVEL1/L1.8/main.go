package main

import "fmt"

/*
Дана переменная типа int64. Разработать программу, которая устанавливает i-й бит этого числа в 1 или 0.
Пример: для числа 5 (0101₂) установка 1-го бита в 0 даст 4 (0100₂).
Подсказка: используйте битовые операции (|, &^).
*/

func main() {
	var x int64
	var i int
	for { //выход из цикла по Ctrl+C
		fmt.Println("Input the number you want to change:")
		fmt.Scan(&x)
		fmt.Println("Which bit(index) should be inversed? Enter only positive number less than 64:")
		fmt.Scan(&i)
		if i > 63 || i < 0 {
			fmt.Println("Invalid index. Try again.")
			continue
		}
		fmt.Printf("Original number is: %d\n%064b\n", x, x)
		x ^= 1 << i
		fmt.Printf("Resulting number is: %d\n%064b\n\n\n", x, x)
	}
}
