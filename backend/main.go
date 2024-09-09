package main

import (
	"fmt"
	"link-shortener-backend/src/handlers"
	"link-shortener-backend/src/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	router := gin.Default()

	router.POST("/api/auth/login", handlers.Login)
	router.POST("/api/auth/register", handlers.Register)

	privateGroup := router.Group("/api/")
	privateGroup.Use(handlers.AuthMiddleware())
	privateGroup.POST("/links/create", handlers.CreateLink)
	privateGroup.POST("/clicks/create", handlers.CreateClick)
	privateGroup.GET("/links/get/:id", handlers.GetLink)
	privateGroup.GET("/links/all", handlers.GetAllLinks)

	privateGroup.POST("/redirects/create", handlers.CreateRedirect)
	privateGroup.GET("/redirects/get/:linkID", handlers.GetRedirectsByLinkID)
	privateGroup.DELETE("/redirects/delete/:redirectID", handlers.DeleteRedirect)
	privateGroup.PUT("/redirects/update/:redirectID", handlers.UpdateRedirect)

	privateGroup.POST("/analytics/get", handlers.GetStatistics)
	privateGroup.POST("/analytics/daily", handlers.GetDailyStatistics)
	// This handles everything related to the shortened link
	router.GET("/:shortId", handlers.Redirect)

	repository.InitDatabase()
	router.Run(":8080")
}
