// internal/service/service_test.go
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"orderservice/internal/cache"
	"orderservice/internal/model"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

// простой фейк под интерфейс репозитория
type fakeRepo struct {
	AddNewOrderFunc         func(ctx context.Context, o *model.Order) error
	GetOrderInfoFunc        func(ctx context.Context, uid string) (*model.Order, error)
	GetAllOrdersFunc        func(ctx context.Context) ([]model.Order, error)
	PushOrderToRawTableFunc func(ctx context.Context, brokenOrder model.InvalidRequest) error
}

func (f *fakeRepo) AddNewOrder(ctx context.Context, o *model.Order) error {
	if f.AddNewOrderFunc != nil {
		return f.AddNewOrderFunc(ctx, o)
	}
	return nil
}
func (f *fakeRepo) GetOrderByUID(ctx context.Context, uid string) (*model.Order, error) {
	if f.GetOrderInfoFunc != nil {
		return f.GetOrderInfoFunc(ctx, uid)
	}
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeRepo) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	if f.GetAllOrdersFunc != nil {
		return f.GetAllOrdersFunc(ctx)
	}
	return nil, nil
}
func (f *fakeRepo) PushOrderToRawTable(ctx context.Context, broken model.InvalidRequest) error {
	if f.PushOrderToRawTableFunc != nil {
		return f.PushOrderToRawTableFunc(ctx, broken)
	}
	return nil
}

func TestProcessKafkaMessage_OK(t *testing.T) {
	repo := &fakeRepo{
		AddNewOrderFunc: func(ctx context.Context, o *model.Order) error {
			if o.OrderUID == "" {
				t.Fatalf("expected order decoded")
			}
			return nil
		}, GetAllOrdersFunc: func(ctx context.Context) ([]model.Order, error) {
			return nil, nil
		},
	}
	mapa := cache.OrderMap{
		CacheMap: make(map[string]model.Order),
		Repo:     repo,
	}
	svc := NewOrderService(repo, &mapa)
	msg := kafka.Message{
		Value: []byte(`{"order_uid":"u1","track_number":"T","entry":"WBIL","delivery":{"name":"A","phone":"1","zip":"1","city":"C","address":"A","region":"R","email":"e@e"},"payment":{"transaction":"u1","request_id":"","currency":"USD","provider":"p","amount":1,"payment_dt":1637907727,"bank":"b","delivery_cost":1,"goods_total":1,"custom_fee":500},"items":[{"chrt_id":1,"track_number":"T","price":1,"rid":"r","name":"n","sale":0,"size":"s","total_price":1,"nm_id":1,"brand":"b","status":1}],"locale":"en","internal_signature":"","customer_id":"c","delivery_service":"d","shardkey":"1","sm_id":1,"date_created":"2021-11-26T06:22:19Z","oof_shard":"1"}`),
	}
	var testOrder model.Order
	if err := json.Unmarshal(msg.Value, &testOrder); err != nil {
	}
	rawTestOrder, _ := json.Marshal(testOrder)

	svc.AddNewOrder(&msg)
	svcOrder, ok := mapa.CacheMap["u1"]
	if !ok {
		t.Fatalf("expected order created and in cache")
	}

	rawSvcOrder, _ := json.Marshal(svcOrder)
	if string(rawTestOrder) != string(rawSvcOrder) {
		fmt.Println("Original json:", string(rawTestOrder))
		fmt.Println("Processed one:", string(rawSvcOrder))
		t.Fatalf("expected input and output order data to be equal")
	}
}
