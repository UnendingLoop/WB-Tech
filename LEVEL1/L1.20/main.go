package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

/*
Разработать программу, которая переворачивает порядок слов в строке.
Пример: входная строка:
«snow dog sun», выход: «sun dog snow».
Считайте, что слова разделяются одиночным пробелом.
Постарайтесь не использовать дополнительные срезы, а выполнять операцию «на месте».
*/

func main() {
	fmt.Println("Input the words that need to be reverted(more than 1 using space as a separator):")
	str, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Printf("Failed to read the input: %s\nTry relaunching the app.", err)
		return
	}

	if str == "" {
		log.Printf("The input is empty. Try again.")
		return
	}

	str = strings.TrimSpace(str)
	input := strings.Split(str, " ")

	switch len(str) {
	case 0:
		fmt.Println("0 words are detected. Try again.")
		return
	case 1:
		fmt.Println("Only 1 word is detected - try again.")
		return
	}

	for i, j := 0, len(input)-1; i <= j; {
		input[i], input[j] = input[j], input[i]
		i++
		j--
	}
	fmt.Println("The reverted order of words is: ", input)
}
