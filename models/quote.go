package models

import "time"

type Quote struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SymbolID  uint      `gorm:"not null;index" json:"symbol_id"`
	Timestamp time.Time `gorm:"not null;index" json:"timestamp"`

	// Current price
	CurrentPrice float64 `gorm:"not null" json:"c"`

	// Change from previous close
	Change float64 `gorm:"not null" json:"d"`

	// Percent change from previous close
	PercentChange float64 `gorm:"not null" json:"dp"`

	// High price of the day
	HighPrice float64 `gorm:"not null" json:"h"`

	// Low price of the day
	LowPrice float64 `gorm:"not null" json:"l"`

	// Open price of the day
	OpenPrice float64 `gorm:"not null" json:"o"`

	// Previous close price
	PreviousClose float64 `gorm:"not null" json:"pc"`

	// Relationship
	Symbol Symbol `gorm:"foreignKey:SymbolID" json:"symbol,omitempty"`
}
