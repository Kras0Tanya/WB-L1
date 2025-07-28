package main

import (
	"fmt"
	"sync"
	"time"
)

// в коде реализованы два варианта безопасной записи в map с использованием sync.Mutex и sync.Map; проверен на гонки с -race (go run main7.go -race)

// safeMap - структура для безопасной работы с мапой с использованием sync.Mutex
type SafeMap struct {
	mu   sync.Mutex
	data map[int]int // ключ: номер воркера, значение: счетчик записей
}

// newSafeMap создает новый экземпляр SafeMap
func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[int]int),
	}
}

// store записывает значение в мапу с блокировкой
func (sm *SafeMap) Store(key, value int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
}

// load возвращает значение по ключу с блокировкой
func (sm *SafeMap) Load(key int) (int, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	value, exists := sm.data[key]
	return value, exists
}

// print выводит содержимое мапы
func (sm *SafeMap) Print() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	fmt.Println("SafeMap содержимое:", sm.data)
}

// workerMutex - воркер, записывающий в SafeMap
func workerMutex(id int, sm *SafeMap, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		sm.Store(id, i+1) // записываем счетчик для воркера
		// задержка для имитации работы и увеличения вероятности гонок
		time.Sleep(time.Microsecond)
	}
}

// workerSyncMap - воркер, записывающий в sync.Map
func workerSyncMap(id int, sm *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		sm.Store(id, i+1) // Записываем счетчик для воркера
		time.Sleep(time.Microsecond)
	}
}

// printSyncMap выводит содержимое sync.Map
func printSyncMap(sm *sync.Map) {
	fmt.Println("sync.Map содержимое:")
	sm.Range(func(key, value interface{}) bool {
		fmt.Printf("Ключ: %v, Значение: %v\n", key, value)
		return true
	})
}

func main() {
	const numWorkers = 10
	var wg sync.WaitGroup

	// вывод SafeMap с sync.Mutex
	fmt.Println("= SafeMap с sync.Mutex =")
	safeMap := NewSafeMap()
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go workerMutex(i, safeMap, &wg)
	}
	wg.Wait()
	safeMap.Print()
	fmt.Println()

	// вывод sync.Map
	fmt.Println("= sync.Map =")
	syncMap := &sync.Map{}
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go workerSyncMap(i, syncMap, &wg)
	}
	wg.Wait()
	printSyncMap(syncMap)
}

/*
Конкурентная запись в map

Реализовать безопасную для конкуренции запись данных в структуру map.

Подсказка: необходимость использования синхронизации (например, sync.Mutex или встроенная concurrent-map).

Проверьте работу кода на гонки (util go run -race).
*/
