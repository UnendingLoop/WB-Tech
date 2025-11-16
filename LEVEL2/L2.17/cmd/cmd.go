// Package cmd is an entry-point to the main logic of application
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"telnet/config"
	"telnet/start"
)

func StartTelnet() {
	telnetConfig := config.Connection{}

	start.ReadArgs(&telnetConfig)

	conn := telnetConfig.GetConnection()
	defer conn.Close()

	fmt.Printf("Connection established to: %q at %v\n", conn.RemoteAddr(), time.Now().UTC())

	// Канал для Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	// отслеживатель контекста - закрывает соединение если контекст отменен
	go func() {
		<-ctx.Done()
		log.Println("Context cancelled, closing connection...")
		conn.Close()
	}()

	// запуск слушателя прерываний:
	wg.Add(1)
	go interruptWarden(sigCh, ctx, &wg, cancel)

	// 2 гоуртины для STDIN и STDOUT:
	wg.Add(1)
	go localScanner(conn, ctx, &wg, cancel)
	wg.Add(1)
	go netReader(conn, ctx, &wg, cancel)

	wg.Wait()
	log.Print("All goroutines are stopped. Exiting app...")
}

func localScanner(c net.Conn, ctx context.Context, wg *sync.WaitGroup, cancel context.CancelFunc) {
	defer wg.Done()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			// log.Println("Context cancelled. Exiting keyboard scanner...")
			return
		default:
			if !scanner.Scan() { // ловим Ctrl+D (Сtrl+Z на винде)
				// log.Println("End of input. Cancelling context...")
				cancel()
				return
			}
			line := strings.TrimSpace(scanner.Text())

			if len(line) > 0 {
				_, err := c.Write([]byte(line + "\r\n"))
				if err != nil {
					select {
					case <-ctx.Done():
						return
					default:
						if err == io.EOF {
							log.Println("Connection closed by remote host")
							cancel()
							return
						}
						if !strings.Contains(err.Error(), "use of closed network connection") {
							log.Printf("Read error: %v", err)
						}
						return
					}
				}
			}
		}
	}
}

func netReader(c net.Conn, ctx context.Context, wg *sync.WaitGroup, cancel context.CancelFunc) {
	defer wg.Done()

	reader := bufio.NewReader(c)

	for {
		select {
		case <-ctx.Done():
			// log.Println("Context cancelled. Exiting netreader...")
			return
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				select {
				case <-ctx.Done():
					// Контекст уже отменён - просто выходим
					return
				default:
					if err == io.EOF {
						log.Println("Connection closed by remote host")
						cancel()
						return
					}
					if !strings.Contains(err.Error(), "use of closed network connection") {
						log.Printf("Read error: %v", err)
					}
					return
				}
			}

			fmt.Println(">> ", line)
		}
	}
}

func interruptWarden(s chan os.Signal, ctx context.Context, wg *sync.WaitGroup, cancel context.CancelFunc) {
	defer wg.Done()

	for {
		select {
		case <-s:
			// fmt.Println("Interruption received. Cancelling context...")
			cancel()
			return
		case <-ctx.Done():
			// log.Println("Context cancelled. Exiting interruption listener...")
			return
		}
	}
}
