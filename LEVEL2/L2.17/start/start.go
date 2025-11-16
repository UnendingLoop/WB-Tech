// Package start provides initialization of values for opening connection from os.Args flags
package start

import (
	"flag"
	"log"
	"os"

	"telnet/config"
)

func ReadArgs(telnetConfig *config.Connection) {
	parser := flag.NewFlagSet("Simplified TelNet", flag.ExitOnError)

	parser.StringVar(&telnetConfig.Host, "host", "", "--host <host> - mandatory, устаналивает хост для подключения")
	parser.StringVar(&telnetConfig.Port, "port", "25", "--port <port> - optional, устанавливает порт для подключения; по умолчанию = '25'")
	parser.DurationVar(&telnetConfig.Timeout, "timeout", 0, "--timeout <Ns> - optional, устанавливает таймаут для закрытия соединения")

	if err := parser.Parse(os.Args[1:]); err != nil {
		log.Fatalf("Failed to parse os.Args: %v", err)
	}

	if parser.NArg() != 0 {
		log.Fatal("Usage: telnet [options]")
	}

	if telnetConfig.Host == "" {
		log.Fatal("Warning: Host is a mandatory option!")
	}
	if telnetConfig.Timeout < 0 {
		log.Fatal("Warning: Timeout cannot be a negative number!")
	}
}
