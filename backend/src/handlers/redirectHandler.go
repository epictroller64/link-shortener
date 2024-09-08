package handlers

import (
	"fmt"
	"link-shortener-backend/src/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRedirect(c *gin.Context) {
	body := repository.Redirect{}
	c.BindJSON(&body)
	redirect, err := repository.CreateRedirect(body)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, redirect)
}

func GetRedirectsByLinkID(c *gin.Context) {
	redirects, err := repository.GetRedirectsByLinkID(c.Param("linkID"))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, redirects)
}

func DeleteRedirect(c *gin.Context) {
	err := repository.DeleteRedirect(c.Param("redirectID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Redirect deleted"})
}

func UpdateRedirect(c *gin.Context) {
	body := repository.Redirect{}
	c.BindJSON(&body)
	redirect, err := repository.UpdateRedirect(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, redirect)
}
