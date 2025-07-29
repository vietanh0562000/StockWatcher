package database

import (
	"stockwatcher/database/repository"
	"stockwatcher/models"
	"sync"

	"gorm.io/gorm"
)

var Instance *Storage
var once sync.Once

type Storage struct {
	StockRepo      StockRepository
	StockPriceRepo StockPriceRepository
}

func GetStorageInstance(db *gorm.DB) *Storage {
	once.Do(func() {
		Instance = &Storage{
			StockRepo:      repository.NewPostgresStockRepository(db),
			StockPriceRepo: repository.NewPostgresStockPriceRepository(db),
		}
	})

	return Instance
}

type StockRepository interface {
	CreateStock(stock *models.Stock) error
	GetStockByID(id uint) (*models.Stock, error)
	GetAllStocks() ([]models.Stock, error)
	UpdateStock(stock *models.Stock) error
	DeleteStock(id uint) error
}

type StockPriceRepository interface {
	CreateStockPrice(price *models.StockPrice) error
	GetStockPriceByID(id uint) (*models.StockPrice, error)
	GetPricesByStockID(stockID uint) ([]models.StockPrice, error)
	UpdateStockPrice(price *models.StockPrice) error
	DeleteStockPrice(id uint) error
}
