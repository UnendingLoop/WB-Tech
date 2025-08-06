package kafka

import (
	"orderservice/internal/model"
	"sync"
)

type OrderMap struct {
	Check        map[string]model.Order
	sync.RWMutex //встраиваем методы мютекса для защиты
}

type Kafka struct{}

// CreateOrderMap returns a new map for using as a cache
func CreateOrderMap() *OrderMap {
	simpleMap := make(map[string]model.Order)
	return &OrderMap{Check: simpleMap}
}
