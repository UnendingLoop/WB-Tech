package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*Разработать программу, которая будет последовательно отправлять значения в канал,
а с другой стороны канала – читать эти значения. По истечении N секунд программа должна завершаться.
Подсказка: используйте time.After или таймер для ограничения времени работы.*/

func main() {
	wg := sync.WaitGroup{}
	channel := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	wg.Add(1)
	go producer(ctx, channel, &wg)
	wg.Add(1)
	go consumer(ctx, channel, &wg)

	wg.Wait()
	fmt.Println("Exiting the program. Bye!")
}

// Producer
func producer(ctx context.Context, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(ch)
	counter := 0
	for {
		counter++
		select {
		case <-ctx.Done():
			fmt.Println("Received timeout. Exiting Producer...")
			return
		case ch <- counter:
			fmt.Println("Sent to channel:", counter)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// Consumer
func consumer(ctx context.Context, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Received timeout. Exiting Consumer...")
			return
		case i := <-ch:
			fmt.Println("Received from channel:", i)
		}
	}
}
