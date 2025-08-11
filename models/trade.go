package models

import (
	"time"
)

type Trade struct {
	ID              uint      `gorm:"primaryKey" json:"i"`
	TradeConditions []string  `gorm:"serializer:json" json:"c"`
	Price           float64   `gorm:"not null" json:"p"`
	Size            float64   `gorm:"not null" json:"s"`
	Timestamp       time.Time `gorm:"not null" json:"t"`
	ExchangeCode    string    `gorm:"not null" json:"x"`
	Tape            string    `gorm:"not null" json:"z"`
	SymbolName      string    `gorm:"not null" json:"symbol_name"`

	Symbol Symbol `gorm:"foreignKey:SymbolName; references:Symbol" json:"stock,omitempty"`
}

func (Trade) TableName() string {
	return "trades"
}
