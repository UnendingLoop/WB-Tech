package main

import "strings"

/*
Рассмотреть следующий код и ответить на вопросы: к каким негативным последствиям он может привести и как это исправить?

Приведите корректный пример реализации.

var justString string

func someFunc() {
  v := createHugeString(1 &lt;&lt; 10)
  justString = v[:100]
}

func main() {
  someFunc()
}
Вопрос: что происходит с переменной justString?
*/

var justString string

func someFunc() {
	v := createHugeString(1 << 7)        // 7й степени двойки достаточно: 128 сэкономит память в createHugeString
	justString = string([]rune(v)[:100]) // скопирует(!) 100 рун, а не байт - исключит битые символы в строке(если важно для логики), если в createHugeString используется UTF-8 не 1-байтный символ
}

func createHugeString(size int) string {
	return strings.Repeat("a", size)
}

func main() {
	someFunc()
}
