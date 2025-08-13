package main

import (
	"fmt"
	"reflect"
)

// использование переключателя типов с детекцией каналов через reflect
func GetTypeUsingSwitch(value interface{}) string {
	switch value.(type) {
	case int:
		return "int"
	case string:
		return "string"
	case bool:
		return "bool"
	default:
		t := reflect.TypeOf(value)
		if t != nil && t.Kind() == reflect.Chan {
			return "chan " + t.Elem().String()
		}
		return "unknown"
	}
}

// использование рефлексии (с указанием содержимого канала)
func GetTypeUsingReflect(value interface{}) string {
	if value == nil {
		return "unknown"
	}
	t := reflect.TypeOf(value)
	if t == nil {
		return "unknown"
	}
	if t.Kind() == reflect.Chan {
		return "chan " + t.Elem().String()
	}
	switch t.Kind() {
	case reflect.Int:
		return "int"
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "bool"
	default:
		return "unknown"
	}
}

// пользовательский тип для теста
type MyInt int

func main() {
	// создаём каналы разных типов
	sendOnlyChan := make(chan<- int)            // только отправка
	recvOnlyChan := make(<-chan int)            // только приём
	structChan := make(chan struct{})           // канал со структурой
	pointerChan := make(chan *int)              // канал с указателем
	var myIntChan chan MyInt = make(chan MyInt) // канал с пользовательским типом
	bufferedChan := make(chan int, 10)          // буферизированный канал
	chanOfChan := make(chan chan int)           // канал с каналом
	sliceChan := make(chan []int)               // канал со слайсом
	mapChan := make(chan map[string]int)        // канал с мапой
	funcChan := make(chan func())               // канал с функцией

	slice := []interface{}{
		123,                // int
		"Hello, Go!",       // string
		false,              // bool
		make(chan int),     // chan int
		make(chan string),  // chan string
		make(chan bool),    // chan bool
		1.1,                // float64 (unknown)
		make(chan float64), // chan float64
		nil,                // nil (unknown)
		sendOnlyChan,       // chan<- int
		recvOnlyChan,       // <-chan int
		structChan,         // chan struct{}
		pointerChan,        // chan *int
		myIntChan,          // chan MyInt
		bufferedChan,       // chan int (буферизированный)
		chanOfChan,         // chan chan int
		sliceChan,          // chan []int
		mapChan,            // chan map[string]int
		funcChan,           // chan func()
	}

	fmt.Println("Using type switch:")
	for _, value := range slice {
		fmt.Printf("%v - %s\n", value, GetTypeUsingSwitch(value))
	}

	fmt.Println("\nUsing reflect:")
	for _, value := range slice {
		fmt.Printf("%v - %s\n", value, GetTypeUsingReflect(value))
	}
}

/*
Определение типа переменной в runtime
Разработать программу, которая в runtime способна определить тип переменной, переданной в неё (на вход подаётся interface{}).
Типы, которые нужно распознавать: int, string, bool, chan (канал).
Подсказка: оператор типа switch v.(type) поможет в решении.
*/
