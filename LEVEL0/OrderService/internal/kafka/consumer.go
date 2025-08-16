package kafka

import (
	"context"
	"log"
	"orderservice/internal/service"
)

// StartConsumer initializes listening to Kafka messages, which will be forwarded to Service-layer
func StartConsumer(ctx context.Context, srv service.OrderService) {
	reader := NewKafkaReader()
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Kafka read error: %v", err)
			continue
		}
		srv.AddNewOrder(&msg)
		reader.CommitMessages(ctx, msg)
	}
}
