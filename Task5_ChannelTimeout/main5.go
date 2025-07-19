package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

/* в качестве основы взяла своё решение L1.4, поскольку там уже был реализован context;
здесь же добавила таймаут, в соответствии с условиями задачи (см. фабулу ниже кода)
*/

// reader читает значения из канала и выводит их в stdout
// завершается при закрытии канала или отмене контекста
func reader(ch <-chan int, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case data, ok := <-ch:
			if !ok { // если канал закрыт, он завершает работу
				return
			}
			fmt.Printf("Горутина получила значение: %d\n", data)
		case <-ctx.Done(): // контекст отменён (таймаут или ошибка)
			return
		}
	}
}

func main() {

	// проверяем, передан ли аргумент для таймаута
	if len(os.Args) < 2 {
		fmt.Println("Ошибка: укажите количество секунд для таймаута (например, go run main.go 5)")
		os.Exit(1)
	}

	// получаем таймаут из аргумента командной строки
	timeoutSecondsStr := os.Args[1]
	timeoutSeconds, err := strconv.Atoi(timeoutSecondsStr)
	if err != nil {
		fmt.Printf("Ошибка: %s не является числом\n", timeoutSecondsStr)
		os.Exit(1)
	}
	if timeoutSeconds <= 0 {
		fmt.Println("Ошибка: количество секунд должно быть больше 0")
		os.Exit(1)
	}

	// преобразуем таймаут в time.Duration
	timeout := time.Duration(timeoutSeconds) * time.Second

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // гарантируем вызов cancel для освобождения ресурсов

	// создаем канал для данных (в нашем примере - чисел)
	dataCh := make(chan int)

	// создаем WaitGroup для синхронизации
	var wg sync.WaitGroup

	// запускаем ридер
	wg.Add(1)
	go reader(dataCh, ctx, &wg)

	// запускаем горутину для записи чисел в канал
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; ; i++ {
			select {
			case <-ctx.Done(): // контекст отменён (таймаут)
				close(dataCh) // закрываем канал данных
				return
			default:
				dataCh <- i                        // отправляем число в канал
				time.Sleep(500 * time.Millisecond) // задержка для наглядности вывода
			}
		}
	}()

	// ждем завершения всех горутин
	wg.Wait()
	fmt.Println("Программа завершена после таймаута")
}

/*
Таймаут на канал

Разработать программу, которая будет последовательно отправлять значения в канал, а с другой стороны канала – читать эти значения.
По истечении N секунд программа должна завершаться.
Подсказка: используйте time.After или таймер для ограничения времени работы.
*/
