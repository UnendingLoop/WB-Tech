package main

/*Программа должна корректно завершаться по нажатию Ctrl+C (SIGINT).
Выберите и обоснуйте способ завершения работы всех горутин-воркеров при получении сигнала прерывания.
Подсказка: можно использовать контекст (context.Context) или канал для оповещения о завершении.*/
import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	//Парсинг параметра из командной строки
	workers := flag.Int("workers", 5, "number of workers") //запуск: "go run main.go -workers=10"
	flag.Parse()

	channel := make(chan int)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	//Создаем канал для слушания прерываний
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	defer close(sig)
	var shutdownFlag atomic.Bool //для цикла отправки сообщений в канал

	//Запускаем слушатель - при прерывании отменяет контекст
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			<-sig //слушаем Ctrl+C и sigterm
			cancel()
			shutdownFlag.Store(true)
			return
		}
	}()

	//Запуск воркеров
	for i := range *workers {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done(): //при отмене контекста - завершение работы
					fmt.Printf("Worker #%d is shutting down\n", n)
					return
				case i, ok := <-channel:
					if i != 0 { //выполнение кода с прочтением последнего записанного значения, даже если канал закрылся
						fmt.Printf("Worker #%d read from channel value: %d\n", n, i)
					}
					if !ok {
						return //завершение работы, если канал закрыт
					}
				}
			}
		}(i)
	}

	//Запуск бесконечного цикла отправки сообщений в канал
	msg := 1
	for !shutdownFlag.Load() {
		channel <- msg
		time.Sleep(1 * time.Second)
		msg++
	}
	//Закрываем канал после того как прервется цикл отправки сообщений по прерыванию
	close(channel)

	//Ждем все горутины
	wg.Wait()
	fmt.Println("Shutting down the application. Bye!")
}
