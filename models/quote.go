package models

import "time"

type Quote struct {
	ID         uint   `gorm:"primaryKey autoincrement" json:"id"`
	SymbolName string `gorm:"not null;index" json:"symbol_name"`

	AskPrice    float64   `gorm:"not null" json:"ap"`
	AskSize     int64     `gorm:"not null" json:"as"`
	AskExchange string    `gorm:"not null" json:"ax"`
	BidPrice    float64   `gorm:"not null" json:"bp"`
	BidSize     int64     `gorm:"not null" json:"bs"`
	BidExchange string    `gorm:"not null" json:"bx"`
	Conditions  []string  `gorm:"serializer:json" json:"c"`
	Timestamp   time.Time `gorm:"not null;index" json:"t"`
	// Relationship
	Symbol Symbol `gorm:"foreignKey:SymbolName; references:Symbol" json:"symbol,omitempty"`
}
