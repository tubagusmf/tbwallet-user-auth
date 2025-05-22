package database

import (
	"log"

	"github.com/tubagusmf/tbwallet-user-auth/internal/helper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres() *gorm.DB {
	dsn := helper.GetConnectionString()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}
