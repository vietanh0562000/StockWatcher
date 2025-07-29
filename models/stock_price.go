package models

import "time"

type StockPrice struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StockID   uint      `gorm:"not null;index" json:"stock_id"`
	Timestamp time.Time `gorm:"not null;index" json:"timestamp"`
	Open      float64   `gorm:"not null" json:"open"`
	Close     float64   `gorm:"not null" json:"close"`
	High      float64   `gorm:"not null" json:"high"`
	Low       float64   `gorm:"not null" json:"low"`
	Volume    int64     `gorm:"not null" json:"volume"`

	// Relationship
	Stock Stock `gorm:"foreignKey:StockID" json:"stock,omitempty"`
}
