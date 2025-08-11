package repository

import (
	"stockwatcher/models"

	"gorm.io/gorm"
)

type TradeRepository struct {
	db *gorm.DB
}

func NewTradeRepository(db *gorm.DB) *TradeRepository {
	return &TradeRepository{db: db}
}

func (r *TradeRepository) CreateTrade(trade *models.Trade) error {
	return r.db.Create(trade).Error
}

func (r *TradeRepository) GetTradeByID(id uint) (*models.Trade, error) {
	var trade models.Trade
	err := r.db.First(&trade, id).Error
	if err != nil {
		return nil, err
	}
	return &trade, nil
}

func (r *TradeRepository) GetTradesBySymbolID(symbolID uint) ([]models.Trade, error) {
	var trades []models.Trade
	err := r.db.Where("symbol_id = ?", symbolID).Find(&trades).Error
	return trades, err
}

func (r *TradeRepository) UpdateTrade(trade *models.Trade) error {
	return r.db.Save(trade).Error
}

func (r *TradeRepository) DeleteTrade(id uint) error {
	return r.db.Delete(&models.Trade{}, id).Error
}
