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
	SymbolRepo ISymbolRepository
	TradeRepo  ITradeRepository
	QuoteRepo  IQuoteRepository
}

func GetStorageInstance(db *gorm.DB) *Storage {
	once.Do(func() {
		Instance = &Storage{
			SymbolRepo: repository.NewSymbolRepository(db),
			TradeRepo:  repository.NewTradeRepository(db),
			QuoteRepo:  repository.NewQuoteRepository(db),
		}
	})

	return Instance
}

type ISymbolRepository interface {
	CreateSymbol(symbol *models.Symbol) error
	GetSymbolByID(id uint) (*models.Symbol, error)
	GetAllSymbols() ([]models.Symbol, error)
	UpdateSymbol(symbol *models.Symbol) error
	DeleteSymbol(id uint) error
}

type ITradeRepository interface {
	CreateTrade(trade *models.Trade) error
	GetTradeByID(id uint) (*models.Trade, error)
	GetTradesBySymbolID(symbolID uint) ([]models.Trade, error)
	UpdateTrade(trade *models.Trade) error
	DeleteTrade(id uint) error
}

type IQuoteRepository interface {
	CreateQuote(quote *models.Quote) error
	GetQuoteByID(id uint) (*models.Quote, error)
	GetQuotesBySymbolID(symbolID uint) ([]models.Quote, error)
	GetLatestQuoteBySymbolID(symbolID uint) (*models.Quote, error)
	UpdateQuote(quote *models.Quote) error
	DeleteQuote(id uint) error
}
