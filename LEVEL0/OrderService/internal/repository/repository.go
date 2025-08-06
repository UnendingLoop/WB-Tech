package repository

import (
	"context"
	"orderservice/internal/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

// GetOrderByUID finds order by its UUID and provides it with error message(if any)
func (OR *OrderRepository) GetOrderByUID(ctx context.Context, uid string) (*model.Order, error) {
	var order model.Order
	if err := OR.DB.WithContext(ctx).Preload("Delivery").Preload("Payment").Preload("Items").Where("order_uid = ?", uid).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// AddNewOrder creates a new record in DB using ctx
func (OR *OrderRepository) AddNewOrder(ctx context.Context, neworder *model.Order) error {
	tx := OR.DB.WithContext(ctx).Begin()

	if err := tx.Create(&neworder).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(&neworder.Delivery).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(&neworder.Payment).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(&neworder.Items).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// GetAllOrders retreives existing orders from DB with limit=300
func (OR *OrderRepository) GetAllOrders(newmap map[string]model.Order) ([]model.Order, error) {
	var orders []model.Order
	if err := OR.DB.Preload("Delivery").Preload("Payment").Preload("Items").Limit(300).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
