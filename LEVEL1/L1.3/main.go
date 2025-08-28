package main

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

/*Реализовать постоянную запись данных в канал (в главной горутине).
Реализовать набор из N воркеров, которые читают данные из этого канала и выводят их в stdout.
Программа должна принимать параметром количество воркеров и при старте создавать указанное число горутин-воркеров.*/

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
					if ok { //выполнение кода только если канал открыт
						fmt.Printf("Worker #%d read from channel value: %d\n", n, i)
						continue
					}
					return //завершение работы, если канал закрыт
				}
			}
		}(i)
	}

	//Запуск бесконечного цикла отправки сообщений в канал
	msg := 0
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
