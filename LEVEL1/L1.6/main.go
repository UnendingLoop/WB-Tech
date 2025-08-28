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

/*Реализовать все возможные способы остановки выполнения горутины.
Классические подходы: выход по условию, через канал уведомления, через контекст, прекращение работы runtime.Goexit() и др.
Продемонстрируйте каждый способ в отдельном фрагменте кода.*/

func main() {
	wg := sync.WaitGroup{}

	//По каналу уведомления
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	defer close(sig)

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Launching goroutine #1")
		<-sig
		fmt.Println("Received interruption - exiting goroutine #1.\n")
	}()
	wg.Wait() //ждем и нажимаем Ctrl+C, переходим в следующему способу

	//По условию
	var toShutdown atomic.Bool //чтобы избежать гонки

	wg.Add(1)
	go func() {
		fmt.Println("Launching goroutine #2")
		defer wg.Done()
		for !toShutdown.Load() {
			fmt.Println("Gouroutine #2 waiting for FALSE...")
			time.Sleep(1 * time.Second)
		}
		fmt.Println("Gouroutine #2 read FALSE. Shutting down.\n")
	}()

	time.Sleep(3 * time.Second)
	toShutdown.Store(true)
	wg.Wait() //ждем таймер и установку toShutdown в значение true. переходим к следующему куйсу

	//закрытие читаемого канала
	channel := make(chan struct{})

	wg.Add(1)
	go func() {
		fmt.Println("Launching goroutine #3")
		defer wg.Done()
		for {
			_, ok := <-channel
			switch ok {
			case true:
				fmt.Println("Channel is open, goroutine #3 continue working...")
				time.Sleep(1 * time.Second)
			case false:
				fmt.Println("Channel is closed, goroutine #3 shutting down.\n")
				return
			}
		}
	}()
	time.Sleep(3 * time.Second)
	close(channel)
	wg.Wait()

	//контекст с таймаутом и с отменой
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Launching goroutine #4")
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Gouroutine #4 received timeout. Shutting down.\n")
				return
			default:
				fmt.Println("Gouroutine #4 waiting for timeout...")
				time.Sleep(1 * time.Second)
			}
		}
	}()
	time.Sleep(3 * time.Second)
	cancel() //вручную отменяем контекст не дожидаясь таймаута, но для горутины не будет никакой разницы
	wg.Wait()

	//по таймеру
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Launching goroutine #5")
		<-time.After(5 * time.Second)
		fmt.Println("Goroutine #5 received timer timeout. Shutting down.\n")
		return
	}()
	wg.Wait()

	//паника
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Launching goroutine #6")
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Goroutine #6 caught panic: %v\nShutting down...\n\n", r)
				return
			}
		}()
		fmt.Println("Goroutine #6 waiting for panic...")
		time.Sleep(1 * time.Second)
		panic("Time to shutdown!")
	}()
	wg.Wait()

	//runtime.Goexit() - defer выполняются
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer fmt.Println("Goroutine #7 shutting down after Goexit.\n")
		fmt.Println("Launching goroutine #7")
		runtime.Goexit()
		fmt.Println("Goroutine #7: this line will not be printed")
	}()
	wg.Wait()

	//os.Exit(code) - без defer-ов так как они не выполнятся
	go func() {
		fmt.Println("Launching goroutine #8")
		fmt.Println("Goroutine #8 and the whole app will be shut down in 1 sec.")
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
	time.Sleep(2 * time.Second)
}
