package models

import "time"

// Transaction represents a financial transaction in the system.
type Transaction struct {
	ID            uint      `gorm:"primaryKey"`
	FromAccountID *uint     `gorm:"index"` // optional
	ToAccountID   *uint     `gorm:"index"` // optional
	Amount        float64   `gorm:"not null"`
	Type          string    `gorm:"not null"` // e.g., "transfer", "withdraw", "deposit"
	Description   string    `gorm:"type:text"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`

	FromAccount *Account `gorm:"foreignKey:FromAccountID"`
	ToAccount   *Account `gorm:"foreignKey:ToAccountID"`
}
