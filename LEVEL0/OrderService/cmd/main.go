package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	handler "orderservice/internal/api"
	"orderservice/internal/cache"
	"orderservice/internal/db"
	"orderservice/internal/kafka"
	"orderservice/internal/repository"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set in env")
	}
	db := db.ConnectPostgres(dsn)
	repo := repository.NewOrderRepository(db)

	orderMap := cache.CreateOrderCache(repo)
	if err := cache.WarmUpCache(orderMap); err != nil {
		log.Fatalf("Failed to load cache: %v", err)
	}

	ctx := context.Background()
	go kafka.StartConsumer(ctx, orderMap)

	orderRepo := handler.OrderHandler{
		Repo: repo,
		Map:  orderMap,
	}

	r := chi.NewRouter()
	r.Get("/order/{order_uid}", orderRepo.GetOrderInfo)

	fmt.Println("Server running on http://localhost:8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatal(err)
	}
}
