package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

/*
Разработать программу, которая переворачивает подаваемую на вход строку.
Например: при вводе строки «главрыба» вывод должен быть «абырвалг».
Учтите, что символы могут быть в Unicode (русские буквы, emoji и пр.), то есть
просто iterating по байтам может не подойти — нужен срез рун ([]rune).
*/

func main() {
	fmt.Println("Input the text that needs to be reverted:")
	str, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Printf("Failed to read the input: %s\nTry relaunching the app.", err)
		return
	}
	if str == "" {
		log.Printf("The input is empty. Try again.")
		return
	}
	runes := []rune(str)

	for i, j := 0, len(runes)-1; i <= j; {
		runes[i], runes[j] = runes[j], runes[i]
		i++
		j--
	}
	answer := string(runes)
	fmt.Println("The reverted string is:", answer)
}
