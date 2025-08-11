package repository

import (
	"stockwatcher/models"

	"gorm.io/gorm"
)

type SymbolRepository struct {
	db *gorm.DB
}

func NewSymbolRepository(db *gorm.DB) *SymbolRepository {
	return &SymbolRepository{db: db}
}

func (r *SymbolRepository) CreateSymbol(symbol *models.Symbol) error {
	if symbol.Figi == "" {
		symbol.Figi = symbol.Exchange + symbol.Symbol
	}
	return r.db.Create(symbol).Error
}

func (r *SymbolRepository) GetSymbolByID(id uint) (*models.Symbol, error) {
	var symbol models.Symbol
	err := r.db.First(&symbol, id).Error
	if err != nil {
		return nil, err
	}
	return &symbol, nil
}

func (r *SymbolRepository) GetAllSymbols() ([]models.Symbol, error) {
	var symbols []models.Symbol
	err := r.db.Find(&symbols).Error
	return symbols, err
}

func (r *SymbolRepository) UpdateSymbol(symbol *models.Symbol) error {
	if symbol.Figi == "" {
		symbol.Figi = symbol.Exchange + symbol.Symbol
	}
	return r.db.Save(symbol).Error
}

func (r *SymbolRepository) DeleteSymbol(id uint) error {
	return r.db.Delete(&models.Symbol{}, id).Error
}
