package kafka

import (
	"context"
	"encoding/json"
	"log"
	"orderservice/internal/model"

	"github.com/segmentio/kafka-go"
)

// PublishOrder used to emulate sending messages from Kafka
func PublishOrder(ctx context.Context, writer *kafka.Writer, order model.Order) error {
	msgBytes, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = writer.WriteMessages(ctx, kafka.Message{
		Value: msgBytes,
	})
	if err != nil {
		log.Printf("Failed to publish order: %v", err)
		return err
	}

	log.Println("Order published to Kafka")
	return nil
}
