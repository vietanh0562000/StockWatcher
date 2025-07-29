package models

import "time"

type Stock struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Symbol    string    `gorm:"uniqueIndex;not null" json:"symbol"`
	Name      string    `gorm:"not null" json:"name"`
	Exchange  string    `gorm:"not null" json:"exchange"`
	Sector    string    `json:"sector"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
