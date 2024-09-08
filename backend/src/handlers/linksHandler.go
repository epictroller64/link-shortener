package handlers

import (
	"fmt"
	"link-shortener-backend/src/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TODO: Overwrite the link creation details with server info
// CreateLink creates a new link
func CreateLink(c *gin.Context) {
	body := repository.Link{}
	c.BindJSON(&body)
	body.CreatedAt = time.Now()
	body.CreatedBy = "test"
	body.Short = "https://short.com/"
	link, err := repository.CreateLink(body)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, link)
}

func CreateClick(c *gin.Context) {
	body := repository.Click{}
	c.BindJSON(&body)
	click, err := repository.CreateClick(body)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	repository.UpdateLinkClickCount(click.LinkID)
	c.JSON(http.StatusOK, click)
}

func GetLink(c *gin.Context) {
	link, err := repository.GetLink(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, link)
}

func GetAllLinks(c *gin.Context) {
	links, err := repository.GetAllLinks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(links)
	c.JSON(http.StatusOK, links)
}
