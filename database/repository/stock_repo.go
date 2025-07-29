package repository

import (
	"stockwatcher/models"

	"gorm.io/gorm"
)

type PostgresStockRepository struct {
	db *gorm.DB
}

func NewPostgresStockRepository(db *gorm.DB) *PostgresStockRepository {
	return &PostgresStockRepository{db: db}
}

func (r *PostgresStockRepository) CreateStock(stock *models.Stock) error {
	return r.db.Create(stock).Error
}

func (r *PostgresStockRepository) GetStockByID(id uint) (*models.Stock, error) {
	var stock models.Stock
	err := r.db.First(&stock, id).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *PostgresStockRepository) GetAllStocks() ([]models.Stock, error) {
	var stocks []models.Stock
	err := r.db.Find(&stocks).Error
	return stocks, err
}

func (r *PostgresStockRepository) UpdateStock(stock *models.Stock) error {
	return r.db.Save(stock).Error
}

func (r *PostgresStockRepository) DeleteStock(id uint) error {
	return r.db.Delete(&models.Stock{}, id).Error
}
