package models

import "time"

type Symbol struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Symbol        string    `gorm:"uniqueIndex;not null" json:"symbol"`
	DisplaySymbol string    `gorm:"not null" json:"displaySymbol"`
	Name          string    `gorm:"not null" json:"description"`
	Currency      string    `gorm:"not null" json:"currency"`
	Exchange      string    `gorm:"not null" json:"mic"`
	Type          string    `gorm:"not null" json:"type"`
	Figi          string    `gorm:"uniqueIndex" json:"figi"`
	Sector        string    `json:"sector"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}
