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
	SymbolRepo     SymbolRepository
	StockPriceRepo StockPriceRepository
	QuoteRepo      QuoteRepository
}

func GetStorageInstance(db *gorm.DB) *Storage {
	once.Do(func() {
		Instance = &Storage{
			SymbolRepo:     repository.NewPostgresSymbolRepository(db),
			StockPriceRepo: repository.NewPostgresStockPriceRepository(db),
			QuoteRepo:      repository.NewPostgresQuoteRepository(db),
		}
	})

	return Instance
}

type SymbolRepository interface {
	CreateSymbol(symbol *models.Symbol) error
	GetSymbolByID(id uint) (*models.Symbol, error)
	GetAllSymbols() ([]models.Symbol, error)
	UpdateSymbol(symbol *models.Symbol) error
	DeleteSymbol(id uint) error
}

type StockPriceRepository interface {
	CreateStockPrice(price *models.StockPrice) error
	GetStockPriceByID(id uint) (*models.StockPrice, error)
	GetPricesByStockID(stockID uint) ([]models.StockPrice, error)
	UpdateStockPrice(price *models.StockPrice) error
	DeleteStockPrice(id uint) error
}

type QuoteRepository interface {
	CreateQuote(quote *models.Quote) error
	GetQuoteByID(id uint) (*models.Quote, error)
	GetQuotesBySymbolID(symbolID uint) ([]models.Quote, error)
	GetLatestQuoteBySymbolID(symbolID uint) (*models.Quote, error)
	UpdateQuote(quote *models.Quote) error
	DeleteQuote(id uint) error
}
