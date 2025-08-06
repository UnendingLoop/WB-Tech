package handler

import (
	"encoding/json"
	"net/http"
	"orderservice/internal/cache"
	"orderservice/internal/repository"

	"github.com/go-chi/chi/v5"
)

// OrderHandler provides a database connection and access to handlers
type OrderHandler struct {
	Repo *repository.OrderRepository
	Map  *cache.OrderMap
}

// GetOrderInfo provides order info by its ID from URL
func (OH *OrderHandler) GetOrderInfo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "order_uid")

	OH.Map.RLock()
	order, ok := OH.Map.Check[id]
	OH.Map.RUnlock()

	if ok {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(order); err != nil {
			http.Error(w, "Failed to encode order info", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	// В кеше нет, идем в бд:
	orderFromDB, err := OH.Repo.GetOrderByUID(r.Context(), id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Обновление кеша
	OH.Map.Lock()
	OH.Map.Check[id] = *orderFromDB
	OH.Map.Unlock()

	// Пишем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderFromDB)
	w.WriteHeader(http.StatusOK)
}
