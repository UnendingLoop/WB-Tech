package main

import (
	"context"
	"fmt"
	"time"
)

/*
Реализовать собственную функцию sleep(duration) аналогично встроенной функции time.Sleep,
которая приостанавливает выполнение текущей горутины.

Важно: в отличии от настоящей time.Sleep, ваша функция должна именно блокировать выполнение
(например, через таймер или цикл), а не просто вызывать time.Sleep :) — это упражнение.

Можно использовать канал + горутину, или цикл на проверку времени (не лучший способ, но для обучения).
*/

func sleepWithContext(timer time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timer)
	defer cancel()
	<-ctx.Done()
}

func sleepWithChan(timer time.Duration) {
	channel := make(chan struct{}, 1)
	go func() {
		ticker := time.NewTimer(timer)
		<-ticker.C
		channel <- struct{}{}
	}()
	<-channel
}

func sleepWithCycle(timer time.Duration) {
	end := time.Now().Add(timer)
	for {
		if time.Now().After(end) || time.Now().Equal(end) {
			time.Sleep(100 * time.Millisecond) // можно конечно обойтись без time.Sleep, но не хочется грузить ЦП
			return
		}
	}
}

func main() {
	fmt.Println("Starting a sleep with context at", time.Now().Truncate(time.Second))
	sleepWithContext(5 * time.Second)
	fmt.Println("Finished sleeping with context at", time.Now().Truncate(time.Second))

	fmt.Println("\nStarting a sleep with channel at", time.Now().Truncate(time.Second))
	sleepWithChan(5 * time.Second)
	fmt.Println("Finished sleeping with channel at", time.Now().Truncate(time.Second))

	fmt.Println("\nStarting a sleep with cycle at", time.Now().Truncate(time.Second))
	sleepWithCycle(5 * time.Second)
	fmt.Println("Finished sleeping with cycle at", time.Now().Truncate(time.Second))
}
