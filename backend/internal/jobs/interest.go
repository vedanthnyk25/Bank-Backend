package jobs

import (
	"backend/internal/constants"
	"backend/internal/db"
	"backend/internal/models"
	"fmt"
	"log"
	"time"
)

func ApplyInterestBatch() {
	const batchSize = 1000
	now := time.Now()
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	dayIndex := int(now.Weekday()) // Sunday = 0
	offset := dayIndex * batchSize

	var accounts []models.Account
	err := db.DB.Where("last_interest_applied_at IS NULL OR last_interest_applied_at < ?", startOfYear).
		Order("id").
		Limit(batchSize).
		Offset(offset).
		Find(&accounts).Error
	if err != nil {
		log.Printf("Error finding accounts: %v", err)
		return
	}

	if len(accounts) == 0 {
		log.Println("No accounts found for interest application.")
		return
	}
	tx := db.DB.Begin()
	for _, acc := range accounts {
		rate := constants.InterestRates[acc.Type]

		interest := acc.Balance * rate
		acc.Balance += interest
		acc.LastInterestAppliedAt = &now

		if err := tx.Save(&acc).Error; err != nil {
			tx.Rollback()
			log.Printf("Error applying interest for account ID %d: %v", acc.ID, err)
			return
		}

		t := models.Transaction{
			ToAccountID: &acc.ID,
			Amount:      interest,
			Type:        "interest",
			Description: fmt.Sprintf("Interest for account ID %d", acc.ID),
		}
		if err := tx.Create(&t).Error; err != nil {
			tx.Rollback()
			log.Printf("Error creating transaction for account ID %d: %v", acc.ID, err)
			return
		}

		tx.Commit()
		log.Printf("Applied interest for account ID %d: %.2f", acc.ID, interest)
	}
}
