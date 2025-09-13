// Package processor iterates through the input string and unpacks letters if necessary e.g.: "ab2c3d\23e"->"abbcccd222e"
package processor

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	errFirstNotLetter   = errors.New("строка начинается с цифры - первым символом должна быть не цифра")
	errTooManySlashes   = errors.New("обнаружено несколько подряд идущих знаков экранирования. Допустим только один обратный слеш, окруженный отличными от него символами")
	errLastRuneIsSlash  = errors.New("строка заканчивается символом экранирования, что недопустимо")
	errLetterAfterSlash = errors.New("обнаружена буква после символа экранирования. После символа экранирования может следовать только цифра")
	errFoundOnlyNumbers = errors.New("строка cодержит только цифры")
	errZeroAsMultiplier = errors.New("обнаружен 0 в качестве множителя - недопустимое значение")
)

func validateString(inputStr string, inputRunes []rune) error {
	// case: only numbers in input string
	if _, err := strconv.Atoi(inputStr); err == nil {
		return errFoundOnlyNumbers
	}
	// case: slash and letter
	if regexp.MustCompile(`\\[A-Za-z]`).MatchString(inputStr) {
		return errLetterAfterSlash
	}
	// case: 0 as multiplyer
	if regexp.MustCompile(`[A-Za-z]0`).MatchString(inputStr) {
		return errZeroAsMultiplier
	}
	// case: start with number
	if unicode.IsDigit(inputRunes[0]) {
		return errFirstNotLetter
	}
	// case: too many slashes
	if strings.Contains(inputStr, "\\\\") {
		return errTooManySlashes
	}
	// case: last symbol is slash
	if inputRunes[len(inputRunes)-1] == '\\' {
		return errLastRuneIsSlash
	}

	return nil
}

// UnZipString - unpacks the input string using the following logic: "a1b2c3d\23e"->"abbcccd222e"
func UnZipString(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	inputRunes := []rune(input)

	// проводим валидацию входной строки
	err := validateString(input, inputRunes)
	if err != nil {
		return "", err
	}

	// готовим буфер для построения результата
	strBuilder := strings.Builder{}

	// итерируемся по рунам
	var currentLetter rune
	var screenIsON bool
	for _, v := range inputRunes {
		// если руна - буква или цифра с включенным жкраниерованием - сразу добавляем в результат и делаем
		// ее currentLetter для последующих повторов, если они есть
		if (!unicode.IsNumber(v) && v != '\\') || screenIsON {
			currentLetter = v
			strBuilder.WriteRune(v)
			screenIsON = false
			continue
		}

		switch v {
		case '\\':
			screenIsON = true
		default: // если дошли до default - значит это точно цифра и нужно делать "умножение" currentLetter
			number, err := strconv.Atoi(string(v))
			if err != nil {
				return "", err
			}
			strBuilder.WriteString(strings.Repeat(string(currentLetter), number-1))
		}
	}

	return strBuilder.String(), nil
}
