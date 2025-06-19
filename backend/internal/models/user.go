package models

import (
	"time"
)

// User represents a user in the system.
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex;not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	Accounts  []Account
	CreatedAt time.Time
}
