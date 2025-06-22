package models

import (
	"time"
)

// Transaction represents a financial transaction in the system.
type Transaction struct {
	ID            uint      `gorm:"primaryKey"`
	FromAccountID uint      `gorm:"not null;index"`
	ToAccountID   uint      `gorm:"not null;index"`
	Amount        float64   `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null"`
}
