package db

import (
	"log"
	"orderservice/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var models = []any{ //the order of tables is important
	&model.Delivery{},
	&model.Payment{},
	&model.Item{},
	&model.Order{},
	&model.InvalidRequest{},
}

// ConnectPostgres creates connection to Postres and runs automigration using structs from order.go
func ConnectPostgres(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatalf("Cannot open db: %v", err)
	}
	if err := db.AutoMigrate(models...); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}
	log.Println("Connected to Postgres")
	return db
}
