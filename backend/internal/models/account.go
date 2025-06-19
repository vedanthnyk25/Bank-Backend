package models

import (
	"time"
)

//Account represents a financial account in the system.

type Account struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"not null"` // Foreign key to User
	Balance   float64 `gorm:"not null"`
	Currency  string  `gorm:"default:'INR';not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
