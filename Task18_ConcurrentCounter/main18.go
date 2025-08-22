package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// вариант 1 - счетчик на основе sync.Mutex
type MutexCounter struct {
	mu    sync.Mutex
	value int
}

// увеличение счётчика (с защитой от гонок с помощью Lock/Unlock)
func (c *MutexCounter) Increment() {
	c.mu.Lock()   // захватываем мьютекс (доступ только одной горутине)
	c.value++     // увеличиваем значение
	c.mu.Unlock() // освобождаем мьютекс
}

// получение текущего значения (тоже с блокировкой)
func (c *MutexCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// вариант 2 -счетчик на основе sync/atomic
type AtomicCounter struct {
	value int64
}

// увеличение счётчика с помощью атомарной операции
func (c *AtomicCounter) Increment() {
	atomic.AddInt64(&c.value, 1)
}

// чтение значения с помощью атомарной загрузки
func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

// демонстрация работы обоих счётчиков
func main() {
	const goroutines = 100  // количество горутин
	const iterations = 1000 // сколько раз каждая горутина увеличит счётчик

	// с MutexCounter
	var wg sync.WaitGroup
	mutexCounter := MutexCounter{}

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				mutexCounter.Increment()
			}
		}()
	}

	wg.Wait()
	fmt.Println("Final value (MutexCounter):", mutexCounter.Value())

	// с AtomicCounter
	atomicCounter := AtomicCounter{}
	wg = sync.WaitGroup{}

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				atomicCounter.Increment()
			}
		}()
	}

	wg.Wait()
	fmt.Println("Final value (AtomicCounter):", atomicCounter.Value())
}

/*
Конкурентный счетчик
Реализовать структуру-счётчик, которая будет инкрементироваться в конкурентной среде (т.е. из нескольких горутин).
По завершению программы структура должна выводить итоговое значение счётчика.
Подсказка: вам понадобится механизм синхронизации, например, sync.Mutex или sync/Atomic для безопасного инкремента.
*/
