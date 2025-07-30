package repository

import (
	"stockwatcher/models"

	"gorm.io/gorm"
)

type PostgresSymbolRepository struct {
	db *gorm.DB
}

func NewPostgresSymbolRepository(db *gorm.DB) *PostgresSymbolRepository {
	return &PostgresSymbolRepository{db: db}
}

func (r *PostgresSymbolRepository) CreateSymbol(symbol *models.Symbol) error {
	if symbol.Figi == "" {
		symbol.Figi = symbol.Exchange + symbol.Symbol
	}
	return r.db.Create(symbol).Error
}

func (r *PostgresSymbolRepository) GetSymbolByID(id uint) (*models.Symbol, error) {
	var symbol models.Symbol
	err := r.db.First(&symbol, id).Error
	if err != nil {
		return nil, err
	}
	return &symbol, nil
}

func (r *PostgresSymbolRepository) GetAllSymbols() ([]models.Symbol, error) {
	var symbols []models.Symbol
	err := r.db.Find(&symbols).Error
	return symbols, err
}

func (r *PostgresSymbolRepository) UpdateSymbol(symbol *models.Symbol) error {
	if symbol.Figi == "" {
		symbol.Figi = symbol.Exchange + symbol.Symbol
	}
	return r.db.Save(symbol).Error
}

func (r *PostgresSymbolRepository) DeleteSymbol(id uint) error {
	return r.db.Delete(&models.Symbol{}, id).Error
}
