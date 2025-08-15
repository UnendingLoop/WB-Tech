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
	"orderservice/internal/service"
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
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to retrieve sql.DB: %v", err)
	}
	defer sqlDB.Close()

	repo := repository.NewOrderRepository(db, dsn)
	orderMap, err := cache.CreateAndWarmUpOrderCache(repo)
	if err != nil {
		log.Fatalf("Failed to load cache: %v", err)
	}

	orderHandler := handler.OrderHandler{
		Service: &service.OrderService{
			Repo: repo,
			Map:  orderMap,
		},
	}

	r := chi.NewRouter()
	r.Get("/order/", orderHandler.GetOrderInfo)

	fmt.Println("Server running on http://localhost:8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatal(err)
	}

	go kafka.StartConsumer(context.Background(), orderHandler.Service)
}
