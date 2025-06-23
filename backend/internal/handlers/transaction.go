package handlers

import (
	"backend/internal/db"
	"backend/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// 1.Transfer funds between accounts
func TransferFunds(c *gin.Context) {
	var input struct {
		FromAccountID uint    `json:"from_account_id"`
		ToAccountID   uint    `json:"to_account_id"`
		Amount        float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || input.FromAccountID == input.ToAccountID || input.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := c.GetUint("user_id")

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var from, to models.Account

		if err := tx.Where("id = ? AND user_id = ?", input.FromAccountID, userID).First(&from).Error; err != nil {
			return fmt.Errorf("Source account not found")
		}
		if err := tx.Where("id = ?", input.ToAccountID).First(&to).Error; err != nil {
			return fmt.Errorf("Destination account not found")
		}

		if from.Balance < input.Amount {
			return fmt.Errorf("Insufficient balance in source account")
		}

		from.Balance -= input.Amount
		to.Balance += input.Amount
		if err := tx.Save(&from).Error; err != nil {
			return fmt.Errorf("Failed to update source account balance")
		}
		if err := tx.Save(&to).Error; err != nil {
			return fmt.Errorf("Failed to update destination account balance")
		}

		transaction := models.Transaction{
			FromAccountID: &input.FromAccountID,
			ToAccountID:   &input.ToAccountID,
			Amount:        input.Amount,
			Type:          "transfer",
			Description:   "Transfer from account " + fmt.Sprint(input.FromAccountID) + " to account " + fmt.Sprint(input.ToAccountID),
		}
		return tx.Create(&transaction).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transfer successful"})
}

// 2. Withdraw funds from an account
func WithDrawFunds(c *gin.Context) {
	var input struct {
		AccountID uint    `json:"account_id"`
		Amount    float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := c.GetUint("user_id")
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var account models.Account
		if err := tx.Where("id = ? AND user_id = ?", input.AccountID, userID).First(&account).Error; err != nil {
			return fmt.Errorf("Account not found")
		}
		if account.Balance < input.Amount {
			return fmt.Errorf("Insufficient balance")
		}

		account.Balance -= input.Amount
		if err := tx.Save(&account).Error; err != nil {
			return fmt.Errorf("Failed to update account balance")
		}
		transaction := models.Transaction{
			FromAccountID: &input.AccountID,
			Amount:        input.Amount,
			Type:          "withdraw",
			Description:   "Withdrawal from account",
		}
		if err := tx.Create(&transaction).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful"})
}

// 3. Deposit funds into an account
func DepositFunds(c *gin.Context) {
	var input struct {
		AccountID uint    `json:"account_id"`
		Amount    float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := c.GetUint("user_id")
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var account models.Account
		if err := tx.Where("id = ? AND user_id = ?", input.AccountID, userID).First(&account).Error; err != nil {
			return fmt.Errorf("Account not found")
		}

		account.Balance += input.Amount
		if err := tx.Save(&account).Error; err != nil {
			return fmt.Errorf("Failed to update account balance")
		}
		transaction := models.Transaction{
			ToAccountID: &input.AccountID,
			Amount:      input.Amount,
			Type:        "deposit",
			Description: "Deposit into account",
		}
		return tx.Create(&transaction).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deposit successful"})
}

// 4. Get transaction history for an account
func GetTransactionHistory(c *gin.Context) {
	accountID := c.Param("id")
	userID := c.GetUint("user_id")

	var account models.Account
	if err := db.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		val, err := strconv.Atoi(l)
		if err == nil && val > 0 {
			limit = val
		}
	}
	if o := c.Query("offset"); o != "" {
		val, err := strconv.Atoi(o)
		if err == nil && val >= 0 {
			offset = val
		}
	}

	var transactions []models.Transaction
	if err := db.DB.Where("from_account_id = ? OR to_account_id = ?", accountID, accountID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
	}
	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

// 5. Get transaction details by ID
func GetTransactionByID(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")

	var transaction models.Transaction
	if err := db.DB.Where("id = ?", id).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	//secure check — only allow if user owns from/to account
	if transaction.FromAccountID != nil {
		var fromAccount models.Account
		if err := db.DB.First(&fromAccount, *transaction.FromAccountID).Error; err == nil && fromAccount.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}
	if transaction.ToAccountID != nil {
		var toAccount models.Account
		if err := db.DB.First(&toAccount, *transaction.ToAccountID).Error; err == nil && toAccount.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}
