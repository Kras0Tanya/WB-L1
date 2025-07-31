package main

import "fmt"

// generateNumbers отправляет числа из массива в канал
func generateNumbers(nums []int, chOut chan<- int) {
	for _, n := range nums {
		chOut <- n
	}
	close(chOut) // закрываем канал после отправки
}

// multiplyByTwo читает числа из канала, умножает на 2 и отправляет в другой канал
func multiplyByTwo(chIn <-chan int, chOut chan<- int) {
	for num := range chIn {
		chOut <- num * 2
	}
	close(chOut) // закрываем канал после обработки
}

func main() {
	numbers := []int{1, 2, 3, 4, 5}

	input := make(chan int)
	output := make(chan int)

	go generateNumbers(numbers, input)
	go multiplyByTwo(input, output)

	// читаем результаты из output и выводим
	for result := range output {
		fmt.Println(result)
	}
}

/*
Конвейер чисел

Разработать конвейер чисел.
Даны два канала: в первый пишутся числа x из массива, во второй – результат операции x*2.
После этого данные из второго канала должны выводиться в stdout.
То есть, организуйте конвейер из двух этапов с горутинами: генерация чисел и их обработка.
Убедитесь, что чтение из второго канала корректно завершается.
*/
