package handlers

import (
	"backend/internal/db"
	"backend/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateAccountInput struct {
	Type     string  `json:"type" binding:"required"`
	Balance  float64 `json:"balance" binding:"required"`
	Currency string  `json:"currency" binding:"required"`
}

// Create account
func CreateAccount(c *gin.Context) {
	var input CreateAccountInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Currency == "" {
		input.Currency = "INR"
	}

	account := models.Account{
		UserID:   c.GetUint("user_id"),
		Type:     input.Type,
		Balance:  input.Balance,
		Currency: input.Currency,
	}

	if err := db.DB.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Account created successfully",
		"account_id": account.ID,
	})
}

// Get accounts
func GetAccounts(c *gin.Context) {
	var accounts []models.Account
	userID := c.GetUint("user_id")

	if err := db.DB.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve accounts"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

// Get account by ID
func GetAccountByID(c *gin.Context) {
	id := c.Param("id")
	var account models.Account
	userID := c.GetUint("user_id")
	if err := db.DB.Where("id = ? AND user_id = ?", id, userID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"account": account})
}
