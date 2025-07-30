package main

import "fmt"

// setBit устанавливает i-й бит числа n в значение value (0 или 1)
func setBit(n int64, i uint, value int) int64 {

	// создаём битовую маску: 1, сдвинутая на i позиций
	mask := int64(1) << i

	if value == 1 {
		return n | mask // установка бита в 1
	} else {
		return n &^ mask // установка бита в 0
	}
}

func main() {
	// реализация примера задачи:
	n := int64(5) // число 5 в двоичной системе: 0101
	i := uint(0)  // позиция бита (1-й бит справа - 0-й по i)
	value := 0    // установить в 0

	result := setBit(n, i, value)
	fmt.Printf("Число %d после установки %d-го бита в %d: %d\n", n, i, value, result)
}
