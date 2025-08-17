package main

import "fmt"

// итеративный бинарный поиск
func binarySearchIterative(arr []int, target int) int {
	left, right := 0, len(arr)-1

	for left <= right {
		mid := left + (right-left)/2
		if arr[mid] == target {
			return mid
		}
		if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}

// рекурсивный бинарный поиск
func binarySearchRecursive(arr []int, target, left, right int) int {
	if left > right {
		return -1
	}

	mid := left + (right-left)/2
	if arr[mid] == target {
		return mid
	}
	if arr[mid] < target {
		return binarySearchRecursive(arr, target, mid+1, right)
	}
	return binarySearchRecursive(arr, target, left, mid-1)
}

func main() {
	arr := []int{1, 3, 5, 7, 9, 11, 13, 15}
	target := 7

	// тестирование итеративного поиска
	resultIter := binarySearchIterative(arr, target)
	fmt.Printf("Итеративный поиск: индекс %d\n", resultIter)

	// тестирование рекурсивного поиска
	resultRec := binarySearchRecursive(arr, target, 0, len(arr)-1)
	fmt.Printf("Рекурсивный поиск: индекс %d\n", resultRec)
}

/*
Бинарный поиск
Реализовать алгоритм бинарного поиска встроенными методами языка.
Функция должна принимать отсортированный слайс и искомый элемент, возвращать индекс элемента или -1, если элемент не найден.
Подсказка: можно реализовать рекурсивно или итеративно, используя цикл for.
*/
