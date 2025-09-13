package main

import (
	"fmt"
)

func main() {
	s := []string{"1", "2", "3"}
	fmt.Println(len(s), cap(s))
	modifySlice(s)
	fmt.Println(s)
}

func modifySlice(i []string) {
	i[0] = "3"
	i = append(i, "4")
	i[1] = "5"
	i = append(i, "6")
}
