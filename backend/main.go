package main

import (
	"link-shortener-backend/src/handlers"
	"link-shortener-backend/src/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	router.POST("/links/create", handlers.CreateLink)
	router.POST("/clicks/create", handlers.CreateClick)
	router.GET("/links/get/:id", handlers.GetLink)
	router.GET("/links/all", handlers.GetAllLinks)

	router.POST("/redirects/create", handlers.CreateRedirect)
	router.GET("/redirects/get/:linkID", handlers.GetRedirectsByLinkID)
	router.DELETE("/redirects/delete/:redirectID", handlers.DeleteRedirect)
	router.PUT("/redirects/update/:redirectID", handlers.UpdateRedirect)

	// This handles everything related to the shortened link
	router.GET("/:shortId", handlers.Redirect)

	repository.InitDatabase()
	router.Run(":8080")
}
