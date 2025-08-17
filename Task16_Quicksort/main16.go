package main

import "fmt"

func quickSort(arr []int) []int {
	// 'базовый случай' (если длина массива 0 или 1, он уже отсортирован)
	if len(arr) <= 1 {
		return arr
	}

	// опорный элемент (берем середину)
	pivot := arr[len(arr)/2]

	// три подмассива (элементы меньше опорного, равные опорному и больше)
	var less, equal, greater []int

	for _, num := range arr {
		switch {
		case num < pivot:
			less = append(less, num)
		case num == pivot:
			equal = append(equal, num)
		case num > pivot:
			greater = append(greater, num)
		}
	}

	// рекурсивная сортировка подмассивов и объединение результатов
	return append(append(quickSort(less), equal...), quickSort(greater)...)
}

func main() {
	arr := []int{9, 7, 5, 11, 12, 2, 14, 3, 10, 6, 5, 5}
	fmt.Println("Исходный массив:", arr)

	sortedArr := quickSort(arr)
	fmt.Println("Отсортированный массив:", sortedArr)
}

/*
Быстрая сортировка (quicksort)
Реализовать алгоритм быстрой сортировки массива встроенными средствами языка. Можно использовать рекурсию.
Подсказка: напишите функцию quickSort([]int) []int которая сортирует срез целых чисел.
Для выбора опорного элемента можно взять середину или первый элемент.
*/
