package kafka

import (
	"context"
	"encoding/json"
	"log"
	"orderservice/internal/cache"
	"orderservice/internal/model"
)

// StartConsumer initializes listening to Kafka messages
func StartConsumer(ctx context.Context, orderMap *cache.OrderMap) {
	reader := NewKafkaReader()
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Kafka read error: %v", err)
			continue
		}

		var order model.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Failed to unmarshal order: %v", err)
			continue
		}
		//Валидация данных; проверка на существование

		// Записываем заказ в базу
		if err := orderMap.Repo.AddNewOrder(ctx, &order); err != nil {
			log.Printf("Failed to save order %s to DB: %v", order.OrderUID, err)
			continue
		} // а если база упала? куда девать прочитанное сообщение? какой таймаут подождать?
		// использовать отдельный кеш чисто для аварийных ситуаций и при восстановлении соединения добавить из этого кега в базу?
		// или писать в файл?

		// Обновление кеша
		orderMap.Lock()
		orderMap.Add[order.OrderUID] = order
		orderMap.Unlock()

		log.Printf("Order saved and cached: %v", order.OrderUID)
	}
}

/*
func validateOrder(order *model.Order) error {

	order.OrderUID != nil
		order.OrderUID == order.Payment.Transaction
		// массив товаров не должен быть пустым
		//у массива товаров одинаковый TrackNumber
		order.TrackNumber == order.Items.TrackNumber
	OrderUID, TrackNumber, Entry, Delivery, Payment, Items
	Name, Phone
	Transaction, Amount, Currency, Provider, PaymentDt
	ChrtID, TrackNumber, Price, Name, TotalPrice

	return nil
}
*/
