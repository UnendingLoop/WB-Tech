package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"main/processor"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	fmt.Println("Введите строку для распаковки:")
	inputStr, err := r.ReadString('\n')
	if err != nil {
		log.Println("Ошибка чтения из stdin:", err)
		return
	}
	result, err1 := processor.UnZipString(inputStr)
	if err1 != nil {
		fmt.Println("Обнаружена ошибка:", err1)
		return
	}
	fmt.Println("Результат распаковки:", result)
}
