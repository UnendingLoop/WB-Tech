package kafka

import (
	"bufio"
	"context"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

// EmulateMsgSending used to emulate real messages flow to test the app in real-time with real DB; mock json-data is read from file
func EmulateMsgSending() {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{kafkaBroker},
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: int(kafka.RequireOne),
		Async:        false,
	})

	file, err := os.Open("./cmd/internal/kafka/mocks.json")
	if err != nil {
		log.Fatalf("Failed to open json-mocks file: %v\nExiting application.", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	counter := 0
	for scanner.Scan() {
		time.Sleep(5 * time.Second)
		counter++
		line := scanner.Bytes()
		err = writer.WriteMessages(context.Background(), kafka.Message{
			Topic: "orders",
			Value: line,
		})
		if err != nil {
			log.Printf("Failed to publish test order #%d: %v", counter, err)
			continue
		}
		log.Printf("Order #%d published to Kafka", counter)
	}

}
