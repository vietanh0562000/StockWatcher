package repository

import (
	"stockwatcher/models"

	"gorm.io/gorm"
)

// --- StockPrice ---

type PostgresStockPriceRepository struct {
	db *gorm.DB
}

func NewPostgresStockPriceRepository(db *gorm.DB) *PostgresStockPriceRepository {
	return &PostgresStockPriceRepository{db: db}
}

func (r *PostgresStockPriceRepository) CreateStockPrice(price *models.StockPrice) error {
	return r.db.Create(price).Error
}

func (r *PostgresStockPriceRepository) GetStockPriceByID(id uint) (*models.StockPrice, error) {
	var price models.StockPrice
	err := r.db.First(&price, id).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

func (r *PostgresStockPriceRepository) GetPricesByStockID(stockID uint) ([]models.StockPrice, error) {
	var prices []models.StockPrice
	err := r.db.Where("stock_id = ?", stockID).Find(&prices).Error
	return prices, err
}

func (r *PostgresStockPriceRepository) UpdateStockPrice(price *models.StockPrice) error {
	return r.db.Save(price).Error
}

func (r *PostgresStockPriceRepository) DeleteStockPrice(id uint) error {
	return r.db.Delete(&models.StockPrice{}, id).Error
}
