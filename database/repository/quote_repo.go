package repository

import (
	"stockwatcher/models"

	"gorm.io/gorm"
)

type QuoteRepository struct {
	db *gorm.DB
}

func NewQuoteRepository(db *gorm.DB) *QuoteRepository {
	return &QuoteRepository{db: db}
}

func (r *QuoteRepository) CreateQuote(quote *models.Quote) error {
	return r.db.Create(quote).Error
}

func (r *QuoteRepository) GetQuoteByID(id uint) (*models.Quote, error) {
	var quote models.Quote
	err := r.db.First(&quote, id).Error
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

func (r *QuoteRepository) GetQuotesBySymbolID(symbolID uint) ([]models.Quote, error) {
	var quotes []models.Quote
	err := r.db.Where("symbol_id = ?", symbolID).Find(&quotes).Error
	return quotes, err
}

func (r *QuoteRepository) GetLatestQuoteBySymbolID(symbolID uint) (*models.Quote, error) {
	var quote models.Quote
	err := r.db.Where("symbol_id = ?", symbolID).Order("timestamp DESC").First(&quote).Error
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

func (r *QuoteRepository) UpdateQuote(quote *models.Quote) error {
	return r.db.Save(quote).Error
}

func (r *QuoteRepository) DeleteQuote(id uint) error {
	return r.db.Delete(&models.Quote{}, id).Error
}
