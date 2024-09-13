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
	router.GET("/api/auth/logout", handlers.Logout)
	router.GET("/api/auth/validate-session", handlers.ValidateSessionHandler)

	// Private routes for authenticated users only
	privateGroup := router.Group("/api/")
	privateGroup.Use(handlers.AuthMiddleware())
	privateGroup.POST("/links/create", handlers.CreateLink)
	privateGroup.POST("/clicks/create", handlers.CreateClick)
	privateGroup.GET("/links/get/:id", handlers.GetLink)
	privateGroup.GET("/links/all", handlers.GetAllLinks)
	privateGroup.DELETE("/links/delete/:id", handlers.DeleteLink)
	privateGroup.GET("/links/recent", handlers.GetRecentLinks)
	privateGroup.POST("/redirects/create", handlers.CreateRedirect)
	privateGroup.GET("/redirects/get/:linkID", handlers.GetRedirectsByLinkID)
	privateGroup.DELETE("/redirects/delete/:redirectID", handlers.DeleteRedirect)
	privateGroup.PUT("/redirects/update/:redirectID", handlers.UpdateRedirect)
	privateGroup.POST("/analytics/get", handlers.GetStatistics)
	privateGroup.POST("/analytics/daily", handlers.GetDailyStatistics)
	privateGroup.POST("/analytics/device", handlers.GetDeviceStatistics)
	privateGroup.GET("/analytics/total", handlers.GetTotalStats)
	// This handles everything related to the shortened link
	privateGroup.POST("/stripe/create-checkout-session", handlers.StripeCreateCheckoutSession)
	privateGroup.GET("/stripe/success", handlers.StripeSuccess)
	privateGroup.GET("/billing/get", handlers.GetBilling)
	privateGroup.GET("/account/get", handlers.GetAccountDetails)
	router.POST("/api/stripe/webhook", handlers.StripeWebHook)
	router.GET("/api/stripe/sync", handlers.StripeSubscriptionSync)
	router.GET("/:shortId", handlers.Redirect)
	router.GET("/api/packages/get", handlers.GetPackages)

	repository.InitDatabase()
	router.Run(":8080")
}
