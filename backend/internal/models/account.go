package models

import "time"

//Account represents a financial account in the system.

type Account struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"not null;index"`
	User      User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Balance   float64 `gorm:"not null"`
	Type      string  `gorm:"not null"`
	Currency  string  `gorm:"default:'INR'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
