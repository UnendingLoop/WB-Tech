package handler

import (
	"context"
	"errors"
	"net/http"
	"orderservice/internal/service"
	"orderservice/internal/web"
)

// OrderHandler provides access to Service layer
type OrderHandler struct {
	Service service.OrderService
}

// GetOrderInfo provides order info by its ID from URL
func (OH *OrderHandler) GetOrderInfo(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	if uid == "" {
		web.Render(w, "search", nil)
		return
	}

	order, err := OH.Service.GetOrderInfo(r.Context(), uid)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			web.Render(w, "error", "Заказ с таким UID не найден")
			return
		case errors.Is(err, context.DeadlineExceeded):
			http.Error(w, err.Error(), http.StatusRequestTimeout)
			return
		default:
			web.Render(w, "error", "Ошибка при поиске заказа: "+err.Error())
			return
		}
	}
	// Успех
	web.Render(w, "order", order)
}
