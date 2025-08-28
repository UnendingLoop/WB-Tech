package calc

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"os"
	"strconv"
	"strings"
)

/*
Разработать программу, которая:
- перемножает,
- делит,
- складывает,
- вычитает
две числовых переменных a, b, значения которых > 2^20 (больше 1 миллиона).

Комментарий: в Go тип int справится с такими числами, но обратите
внимание на возможное переполнение для ещё больших значений.
Для очень больших чисел можно использовать math/big.
*/

// Вычисляем мин и макс значение int для текущей архитектуры
var maxInt = int(^uint(0) >> 1)
var minInt = -maxInt - 1

func add(a, b int) string {
	if !((b > 0 && a > maxInt-b) || (b < 0 && a < minInt-b)) {
		return strconv.Itoa(a + b)
	}

	return bigIntOps(a, b, "+")
}

func subtract(a, b int) string {
	if !((b > 0 && a < minInt+b) || (b < 0 && a > maxInt+b)) {
		return strconv.Itoa(a - b)
	}
	return bigIntOps(a, b, "-")
}

func multiply(a, b int) string {
	if a == 0 || b == 0 {
		return "0"
	}

	mult := a * b
	if mult/a == b {
		return strconv.Itoa(mult)
	}

	return bigIntOps(a, b, "*")

}

func divide(a, b int) string {
	if b == 0 {
		return "cannot divide by 0"
	}
	if a != minInt && b != -1 {
		return strconv.Itoa(a / b)
	}

	return bigIntOps(a, b, "/")
}

func bigIntOps(a, b int, op string) string {
	bigA := big.NewInt(int64(a))
	bigB := big.NewInt(int64(b))
	switch op {
	case "+":
		return big.NewInt(0).Add(bigA, bigB).String()
	case "-":
		return big.NewInt(0).Sub(bigA, bigB).String()
	case "/":
		return big.NewInt(0).Div(bigA, bigB).String()
	case "*":
		return big.NewInt(0).Mul(bigA, bigB).String()
	default:
		return "unknown operation"
	}
}

func input(r io.Reader) (int, int, error) {
	fmt.Println("Input 2 integer numbers for running basic arithmetical operations using space as a separator:")
	input, err := bufio.NewReader(r).ReadString('\n')
	if err != nil {
		return 0, 0, err
	}

	splitInput := strings.Split(strings.TrimSpace(input), " ")

	if len(splitInput) != 2 {
		return 0, 0, fmt.Errorf("Invalid input: expected 2 numbers, got: %d", len(splitInput))
	}

	a, err := strconv.Atoi(splitInput[0])
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to convert input data #1 into integer: %v", err)
	}
	b, err := strconv.Atoi(splitInput[1])
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to convert input data #2 into integer: %v", err)
	}
	return a, b, nil
}

func main() {
	a, b, err := input(os.Stdin)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Multiplication:", multiply(a, b))
	fmt.Println("Division:", divide(a, b))
	fmt.Println("Addition:", add(a, b))
	fmt.Println("Subtraction:", subtract(a, b))
}
