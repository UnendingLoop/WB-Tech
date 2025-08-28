package main

import "fmt"

/*
Разработать программу, которая в runtime способна определить тип переменной,
переданной в неё (на вход подаётся interface{}).

Типы, которые нужно распознавать: int, string, bool, chan (канал).
Подсказка: оператор типа switch v.(type) поможет в решении.
*/

func main() {
	str := "Hello world"
	d := 777
	f := 7.77
	b := true
	channel := make(chan int, 1)
	queue := []any{str, d, f, b, channel, &d}

	for _, v := range queue {
		answer := typeDecoder(v)
		fmt.Println(answer)
	}
}

func typeDecoder(v any) string {
	switch v.(type) {
	case int:
		return "Integer"
	case float64:
		return "Float64"
	case bool:
		return "Boolean"
	case chan int:
		return "Channel for integer"
	case string:
		return "String"
	default:
		return "Undefined"
	}
}
