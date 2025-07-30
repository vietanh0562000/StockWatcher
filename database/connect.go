package database

import (
	"fmt"
	"stockwatcher/config"
	"stockwatcher/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(config *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable",
		config.DatabaseURL, config.DatabaseUser, config.DatabaseName, config.DatabasePort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	// Auto-migrate models
	err = db.AutoMigrate(&models.Symbol{}, &models.StockPrice{}, &models.Quote{})
	if err != nil {
		return nil, err
	}

	fmt.Println("Connect database successful")

	return db, nil
}
