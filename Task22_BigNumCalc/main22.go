package main

import (
	"fmt"
	"math/big"
	"os"
)

// readBigInt читает большое целое число из стандартного ввода и возвращает указатель на big.Int или завершается при ошибке
func readBigInt(prt string) *big.Int {
	fmt.Println(prt)
	n := new(big.Int)
	_, err := fmt.Fscan(os.Stdin, n)
	if err != nil {
		fmt.Printf("Ошибка ввода: %v\n", err)
		os.Exit(1)
	}
	return n
}

func main() {
	// читаем числа
	a := readBigInt("Введите первое число a (должно быть > 2^20):")
	b := readBigInt("Введите второе число b (должно быть > 2^20):")

	// 2^20 = 1_048_576
	limit := new(big.Int).Exp(big.NewInt(2), big.NewInt(20), nil)

	if a.Cmp(limit) <= 0 || b.Cmp(limit) <= 0 {
		fmt.Println("Ошибка: оба числа должны быть строго больше 1 048 576 (2^20)!")
		os.Exit(1)
	}

	fmt.Println("\nРезультаты операций:")
	fmt.Println("a =", a)
	fmt.Println("b =", b)

	fmt.Println("Сложение (a + b)     :", new(big.Int).Add(a, b))
	fmt.Println("Вычитание (a - b)    :", new(big.Int).Sub(a, b))
	fmt.Println("Умножение (a × b)    :", new(big.Int).Mul(a, b))
	fmt.Println("Деление (a ÷ b)      :", new(big.Int).Quo(a, b))
	fmt.Println("Остаток (a mod b)    :", new(big.Int).Rem(a, b))
}

/*
Второй вариант решения, но он не справляется с переполнением.
// используем встроенный int/int64, который работает только пока результат помещается в 64 бита
package main

import (
	"fmt"
	"os"
)

func main() {
	var a, b int64

	fmt.Println("Внимание! Используется int64!")
	fmt.Println("Вероятно переполнение при умножении чисел больше ~10⁹!")
	fmt.Println()

	fmt.Print("Введите a (> 1_048_576): ")
	_, err := fmt.Scan(&a)
	if err != nil || a <= 1_048_576 {
		fmt.Println("Некорректный ввод a")
		os.Exit(1)
	}

	fmt.Print("Введите b (> 1_048_576): ")
	_, err = fmt.Scan(&b)
	if err != nil || b <= 1_048_576 {
		fmt.Println("Некорректный ввод b")
		os.Exit(1)
	}

	fmt.Printf("a = %d\n", a)
	fmt.Printf("b = %d\n", b)

	fmt.Printf("Сложение (a + b)   : %d\n", a+b)
	fmt.Printf("Вычитание (a - b)  : %d\n", a-b)
	fmt.Printf("Умножение (a × b)  : %d   ← может быть НЕВЕРНЫМ при переполнении!\n", a*b)
	fmt.Printf("Деление (a / b)    : %d\n", a/b)
	fmt.Printf("Остаток (a %% b)    : %d\n", a%b)
}
*/

/*
Большие числа и операции
Разработать программу, которая перемножает, делит, складывает, вычитает две числовых переменных a, b,
значения которых > 2^20 (больше 1 миллион).
Комментарий: в Go тип int справится с такими числами, но обратите внимание на возможное переполнение для ещё больших значений.
Для очень больших чисел можно использовать math/big.
*/
