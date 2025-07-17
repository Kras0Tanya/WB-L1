package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

// воркер читает данные из канала с помощью for data := range ch;
// выводит данные в stdout с указанием ID воркера;
// defer wg.Done() уменьшает счетчик WaitGroup после завершения воркера
func worker(id int, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range ch {
		fmt.Printf("Воркер %d обработал число: %d\n", id, data)
	}
}

func main() {

	// проверяем, передан ли аргумент (есть ли аргументы кроме имени программы), иначе - завершаем с кодом ошибки os.Exit(1)
	if len(os.Args) < 2 {
		fmt.Println("Ошибка: укажите количество воркеров (например, go run main.go 5)")
		os.Exit(1)
	}

	// получаем аргумент (количество воркеров): os.Args[1] - строка, содержащая первый аргумент (напр., "5")
	// strconv.Atoi преобразует строку в int, если строка не является числом, возвращается ошибка
	// если err != nil - аргумент некорректен (напр., "ayz"), программа завершается с ошибкой
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

	// создаем канал для данных (в нашем примере - чисел)
	dataCh := make(chan int)

	// создаем WaitGroup для синхронизации
	var wg sync.WaitGroup

	// запускаем воркеры (numWorkers горутин, каждая из которых вызывает функцию worker)
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, dataCh, &wg)
	}

	// буф.канал для сигналов завершения и сами сигналы
	// программа завершается при нажатии Ctrl+C (сигнал SIGINT) или другом сигнале типа kill (SIGTERM).
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// горутина для постоянной записи данных в канал
	go func() {
		for i := 1; ; i++ {
			select {
			case <-sigCh: // если получен сигнал (Ctrl+C)
				close(dataCh) // закрываем канал
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
Работа нескольких воркеров

Реализовать постоянную запись данных в канал (в главной горутине).
Реализовать набор из N воркеров, которые читают данные из этого канала и выводят их в stdout.
Программа должна принимать параметром количество воркеров и при старте создавать указанное число горутин-воркеров.
*/
