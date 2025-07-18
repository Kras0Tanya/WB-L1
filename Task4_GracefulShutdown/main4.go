package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

/* в качестве основы взяла своё решение L1.3, поскольку там уже был реализован graceful shutdown через канал;
здесь же добавила контекст, поскольку это стандартный инструмент, его использование упрощает интеграцию с другими библиотеками,
позволяет добавлять таймауты, дедлайны или значения без изменения структуры программы
*/

// worker запускает воркера с указанным ID, читает данные из канала и обрабатывает их.
// добавлен ctx.Done()
// defer wg.Done() уменьшает счетчик WaitGroup после завершения воркера
func worker(id int, ch <-chan int, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case data, ok := <-ch:
			if !ok { // если канал закрыт, воркер больше не может получать данные, поэтому он завершает работу
				return
			}
			fmt.Printf("Воркер %d обработал число: %d\n", id, data)
		case <-ctx.Done(): // при отмене контекста (например, по Ctrl+C) канал ctx.Done() закрывается, воркер завершает работу
			return
		}
	}
}

func main() {

	// проверяем, передан ли аргумент (есть ли аргументы кроме имени программы), иначе - завершаем с кодом ошибки os.Exit(1)
	if len(os.Args) < 2 {
		fmt.Println("Ошибка: укажите количество воркеров (например, go run main.go 5)")
		os.Exit(1)
	}

	// получаем аргумент (количество воркеров): os.Args[1] - строка, содержащая первый аргумент (напр., "5")
	numWorkersStr := os.Args[1]
	numWorkers, err := strconv.Atoi(numWorkersStr)
	if err != nil {
		fmt.Printf("Ошибка: %s не является числом\n", numWorkersStr)
		os.Exit(1)
	}
	// проверяем, что количество воркеров больше 0
	if numWorkers <= 0 {
		fmt.Println("Ошибка: количество воркеров должно быть больше 0")
		os.Exit(1)
	}

	// Создаем контекст с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())

	// создаем канал для данных (в нашем примере - чисел)
	dataCh := make(chan int)

	// создаем WaitGroup для синхронизации
	var wg sync.WaitGroup

	// запускаем воркеры (numWorkers горутин, каждая из которых вызывает функцию worker)
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, dataCh, ctx, &wg)
	}

	// буф.канал для сигналов завершения и сами сигналы
	// программа завершается при нажатии Ctrl+C (сигнал SIGINT) или другом сигнале типа kill (SIGTERM).
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// горутина для обработки сигналов
	go func() {
		<-sigCh  // ждем сигнала Ctrl+C
		cancel() // отменяем контекст
	}()

	// горутина для постоянной записи данных в канал
	// добавлен ctx.Done()
	go func() {
		for i := 1; ; i++ {
			select {
			case <-ctx.Done(): // если контекст отменён
				close(dataCh) // закрываем канал данных
				return
			default:
				dataCh <- i // записываем данные
			}
		}
	}()

	// ждем завершения воркеров
	wg.Wait()
	fmt.Println("Все воркеры завершили работу!")
}

/*
Завершение по Ctrl+C

Программа должна корректно завершаться по нажатию Ctrl+C (SIGINT).
Выберите и обоснуйте способ завершения работы всех горутин-воркеров при получении сигнала прерывания.
Подсказка: можно использовать контекст (context.Context) или канал для оповещения о завершении.
*/
