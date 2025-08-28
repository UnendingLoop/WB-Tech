package main

import "fmt"

/*
Имеется последовательность строк: ("cat", "cat", "dog", "cat", "tree").
Создать для неё собственное множество.
Ожидается: получить набор уникальных слов.
Для примера, множество = {"cat", "dog", "tree"}.
*/

func main() {
	input := []string{"cat", "cat", "cow", "plant", "cow", "plant", "cow", "fish", "dog", "cat", "tree"}
	uniqueStringsMap := make(map[string]struct{})

	for _, s := range input {
		uniqueStringsMap[s] = struct{}{}
	}

	output := make([]string, 0, len(uniqueStringsMap))

	for k := range uniqueStringsMap {
		output = append(output, k)
	}

	fmt.Printf("Unique strigs from input are: %s\n", output)
}
