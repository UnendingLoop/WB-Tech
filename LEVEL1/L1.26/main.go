package main

import (
	"fmt"
	"unicode"
)

/*
Разработать программу, которая проверяет, что все символы в строке
встречаются один раз (т.е. строка состоит из уникальных символов).

Вывод:
- true, если все символы уникальны,
- false, если есть повторения.

Проверка должна быть регистронезависимой, т.е.
символы в разных регистрах считать одинаковыми, например:
- "abcd" -> true,
- "abCdefAaf" -> false (повторяются a/A),
- "aabcd" -> false.

Подумайте, какой структурой данных удобно воспользоваться для проверки условия.
*/
func isUniqieASCII(input string) bool {
	seen := [128]bool{} // используем просто массив для проверки - быстрее карты и меньше занимает места
	for _, v := range input {
		lowV := unicode.ToLower(v)
		if seen[lowV] {
			return false // как только находим дубликат, сразу выходим не дожидаясь прохождения по строке до конца
		}
		seen[lowV] = true
	}
	return true
}

func isUniqueUnicode(input string) bool {
	seen := make(map[rune]bool) // для Unicode уже придется использовать карту
	for _, v := range input {
		lowV := unicode.ToLower(v) // использую поштучное приведение к lowercase чтобы не генерить доп слайсы при strings.ToLower
		if _, ok := seen[lowV]; ok {
			return false
		}
		seen[lowV] = true
	}
	return true
}

func isUnique(input string) bool {
	for _, v := range input {
		if v > 127 { // проверяем символ: ASCII или Unicode
			return isUniqueUnicode(input)
		}
	}
	return isUniqieASCII(input)
}

func main() {
	lines := []string{"qwerty", "qqwweerrttyy", "Фывапролджэ", "qwertyY", "äåüµ†®éΩ•ü", "ææ@˙~ßß∂¸∞", "123456asdfg", "123sdf321"}
	for _, v := range lines {
		fmt.Printf("Line '%s' uniqueness result: %v\n", v, isUnique(v))
	}
}
