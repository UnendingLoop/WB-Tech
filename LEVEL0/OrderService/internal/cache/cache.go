package cache

import (
	"fmt"
	"orderservice/internal/model"
	"orderservice/internal/repository"
	"sync"
)

// OrderMap provides access to cache-map, contains embedded mutex features
type OrderMap struct {
	Check        map[string]model.Order //Check - более читаемое название при вызове: OH.Map.Check[id]
	Add          map[string]model.Order //Add - более читаемое название при вызове: OH.Map.Add[id]
	Repo         *repository.OrderRepository
	sync.RWMutex //встраиваем методы мютекса для защиты
}

// CreateOrderCache returns a new map for using as a cache, access to DB and embedded mutex
func CreateOrderCache(repo *repository.OrderRepository) *OrderMap {
	simpleMap := make(map[string]model.Order)
	return &OrderMap{Check: simpleMap, Add: simpleMap, Repo: repo}
}

// WarmUpCache loads all records from DB and adds them to cache/map
func WarmUpCache(orderMap *OrderMap) error {
	orders, err := orderMap.Repo.GetAllOrders(orderMap.Check)
	if err != nil {
		return err
	}

	orderMap.Lock()
	for _, v := range orders {
		orderMap.Add[v.OrderUID] = v
	}
	orderMap.Unlock()

	fmt.Println("Cache successfully loaded!")
	return nil
}
