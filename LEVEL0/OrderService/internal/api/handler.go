package handler

import (
	"encoding/json"
	"net/http"
	"orderservice/internal/kafka"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// OrderHandler provides a database connection and access to handlers
type OrderHandler struct { // может обращение в базу вообще тут не надо? перенести в слой кафки? наверное только из мапы надо доставать инфу и отдавать в w
	DB  *gorm.DB //Нужен доступ к БД если нет данных в кэше/мапе
	Map *kafka.OrderMap
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
	//тут нужен запрос в базу через кафку, так как в кеше нет ордера

}
