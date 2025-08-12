package main

import "fmt"

// cоздание множества из слайса строк
func MakeSet(slice []string) []string {
	// cоздание мапы множества из строк
	setMap := make(map[string]struct{})
	for _, value := range slice {
		setMap[value] = struct{}{}
	}
	// перенос элементов множества в слайс
	result := make([]string, 0, len(setMap))
	for key := range setMap {
		result = append(result, key)
	}
	return result
}

func main() {
	slice := []string{"cat", "cat", "dog", "cat", "tree"}
	fmt.Println(MakeSet(slice))
}

/*
Собственное множество строк
Имеется последовательность строк: ("cat", "cat", "dog", "cat", "tree").
Создать для неё собственное множество.
Ожидается: получить набор уникальных слов. Для примера, множество = {"cat", "dog", "tree"}.
*/
