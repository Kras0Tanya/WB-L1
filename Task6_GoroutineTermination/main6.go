package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

/*
В коде реализованы 7 различных способов остановки горутин, включая как классические, так и дополнительные методы.
классические подходы:
- выход по условию (workerCondition с atomic.Bool)
- канал уведомления (workerChannel с каналом done)
- контекст (workerContext с context.Context)
- runtime.Goexit() (workerGoexit)
дополнительные способы:
- закрытие входного канала (workerChannelClose)
- таймер (workerTimer с time.After)
- сигналы ОС (workerSignal)
чаще всего используют контексты или каналы уведомления (считаются идиоматичными и гибкими способами)
*/

// workerCondition завершает горутину по условию (флаг stop с atomic.Bool)
func workerCondition(id int, stop *atomic.Bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; ; i++ {
		if stop.Load() { // атомарно проверяем флаг остановки
			fmt.Printf("Воркер %d (условие): флаг остановки установлен, завершаю работу\n", id)
			return
		}
		fmt.Printf("Воркер %d (условие): обработал число %d\n", id, i)
		time.Sleep(500 * time.Millisecond) // задержка для наглядности
	}
}

// workerChannel завершает горутину при получении сигнала через канал done
func workerChannel(id int, done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; ; i++ {
		select {
		case <-done: // получен сигнал остановки
			fmt.Printf("Воркер %d (канал): получен сигнал остановки, завершаю работу\n", id)
			return
		default:
			fmt.Printf("Воркер %d (канал): обработал число %d\n", id, i)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// workerContext завершает горутину при отмене контекста
func workerContext(id int, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; ; i++ {
		select {
		case <-ctx.Done(): // контекст отменен
			fmt.Printf("Воркер %d (контекст): контекст отменен (%v), завершаю работу\n", id, ctx.Err())
			return
		default:
			fmt.Printf("Воркер %d (контекст): обработал число %d\n", id, i)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// workerGoexit завершает горутину с помощью runtime.Goexit()
func workerGoexit(id int, stop *atomic.Bool, wg *sync.WaitGroup) {
	defer wg.Done() // выполняется даже при Goexit
	for i := 1; ; i++ {
		if stop.Load() { // атомарно проверяем флаг остановки
			fmt.Printf("Воркер %d (Goexit): вызываю Goexit, завершаю работу\n", id)
			runtime.Goexit() // завершает горутину
		}
		fmt.Printf("Воркер %d (Goexit): обработал число %d\n", id, i)
		time.Sleep(500 * time.Millisecond)
	}
}

// workerChannelClose завершает горутину при закрытии входного канала
func workerChannelClose(id int, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range ch { // читаем, пока канал не закрыт
		fmt.Printf("Воркер %d (закрытие канала): обработал число %d\n", id, data)
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Printf("Воркер %d (закрытие канала): канал закрыт, завершаю работу\n", id)
}

// workerTimer завершает горутину по таймауту через time.After
func workerTimer(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	timeout := time.After(2 * time.Second) // таймер на 2 секунды
	for i := 1; ; i++ {
		select {
		case <-timeout: // таймаут истек
			fmt.Printf("Воркер %d (таймер): таймаут истек, завершаю работу\n", id)
			return
		default:
			fmt.Printf("Воркер %d (таймер): обработал число %d\n", id, i)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// workerSignal завершает горутину при получении сигнала ОС
func workerSignal(id int, sigCh <-chan os.Signal, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; ; i++ {
		select {
		case sig := <-sigCh: // получен сигнал ОС
			fmt.Printf("Воркер %d (сигнал): получен сигнал %v, завершаю работу\n", id, sig)
			return
		default:
			fmt.Printf("Воркер %d (сигнал): обработал число %d\n", id, i)
			time.Sleep(100 * time.Millisecond) // уменьшенная задержка для частой проверки
		}
	}
}

func main() {
	var wg sync.WaitGroup

	// демонстрация "выход по условию"
	fmt.Println("= Выход по условию =")
	var stopCondition atomic.Bool
	stopCondition.Store(false)
	wg.Add(1)
	go workerCondition(1, &stopCondition, &wg)
	time.Sleep(2 * time.Second)
	stopCondition.Store(true)
	wg.Wait()
	fmt.Println("Завершено: Выход по условию\n")

	// демонстрация "канал уведомления"
	fmt.Println("= Канал уведомления =")
	done := make(chan struct{})
	wg.Add(1)
	go workerChannel(1, done, &wg)
	time.Sleep(2 * time.Second)
	close(done)
	wg.Wait()
	fmt.Println("Завершено: Канал уведомления\n")

	// демонстрация "контекст"
	fmt.Println("= Контекст =")
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go workerContext(1, ctx, &wg)
	time.Sleep(2 * time.Second)
	cancel()
	wg.Wait()
	fmt.Println("Завершено: Контекст\n")

	// демонстрация "runtime.Goexit()""
	fmt.Println("= runtime.Goexit() =")
	var stopGoexit atomic.Bool
	stopGoexit.Store(false)
	wg.Add(1)
	go workerGoexit(1, &stopGoexit, &wg)
	time.Sleep(2 * time.Second)
	stopGoexit.Store(true)
	wg.Wait()
	fmt.Println("Завершено: runtime.Goexit()\n")

	// демонстрация "закрытие входного канала"
	fmt.Println("= Закрытие входного канала =")
	dataCh := make(chan int)
	wg.Add(1)
	go workerChannelClose(1, dataCh, &wg)
	go func() {
		for i := 1; i <= 5; i++ {
			dataCh <- i
			time.Sleep(500 * time.Millisecond)
		}
		close(dataCh)
	}()
	wg.Wait()
	fmt.Println("Завершено: Закрытие входного канала\n")

	// демонстрация "таймер"
	fmt.Println("= Таймер =")
	wg.Add(1)
	go workerTimer(1, &wg)
	wg.Wait()
	fmt.Println("Завершено: Таймер\n")

	// демонстрация "сигналы ОС"
	fmt.Println("= Сигналы ОС (нажмите Ctrl+C для завершения) =")
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	wg.Add(1)
	go workerSignal(1, sigCh, &wg)
	wg.Wait()
	fmt.Println("Завершено: Сигналы ОС")
}

/*
Остановка горутины

Реализовать все возможные способы остановки выполнения горутины.
Классические подходы: выход по условию, через канал уведомления, через контекст, прекращение работы runtime.Goexit() и др.
Продемонстрируйте каждый способ в отдельном фрагменте кода.
*/
