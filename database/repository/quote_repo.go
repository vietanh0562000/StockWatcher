package repository

import (
	"stockwatcher/models"

	"gorm.io/gorm"
)

type PostgresQuoteRepository struct {
	db *gorm.DB
}

func NewPostgresQuoteRepository(db *gorm.DB) *PostgresQuoteRepository {
	return &PostgresQuoteRepository{db: db}
}

func (r *PostgresQuoteRepository) CreateQuote(quote *models.Quote) error {
	return r.db.Create(quote).Error
}

func (r *PostgresQuoteRepository) GetQuoteByID(id uint) (*models.Quote, error) {
	var quote models.Quote
	err := r.db.First(&quote, id).Error
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

func (r *PostgresQuoteRepository) GetQuotesBySymbolID(symbolID uint) ([]models.Quote, error) {
	var quotes []models.Quote
	err := r.db.Where("symbol_id = ?", symbolID).Find(&quotes).Error
	return quotes, err
}

func (r *PostgresQuoteRepository) GetLatestQuoteBySymbolID(symbolID uint) (*models.Quote, error) {
	var quote models.Quote
	err := r.db.Where("symbol_id = ?", symbolID).Order("timestamp DESC").First(&quote).Error
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

func (r *PostgresQuoteRepository) UpdateQuote(quote *models.Quote) error {
	return r.db.Save(quote).Error
}

func (r *PostgresQuoteRepository) DeleteQuote(id uint) error {
	return r.db.Delete(&models.Quote{}, id).Error
}
