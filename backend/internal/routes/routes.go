package routes

import (
	"backend/internal/handlers"
	"backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// Auth routes
	auth := r.Group("/api/auth")
	auth.Use(middleware.AuthMiddleware())
	auth.GET("/me", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"user_id": c.GetUint("user_id"),
			"role":    c.GetString("role"),
		})
	})

	// Account routes
	account := r.Group("/api/accounts")
	account.Use(middleware.AuthMiddleware())
	account.POST("", handlers.CreateAccount)
	account.GET("", handlers.GetAccounts)
	account.GET("/:id", handlers.GetAccountByID)
}
