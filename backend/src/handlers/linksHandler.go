package handlers

import (
	"fmt"
	"link-shortener-backend/src/repository"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TODO: Overwrite the link creation details with server info
// CreateLink creates a new link
func CreateLink(c *gin.Context) {
	domain := "http://localhost:8080/" // TODO: Change to env variable later or whatever
	shortLink := GenerateShortLink()
	body := repository.Link{}
	c.BindJSON(&body)
	body.CreatedAt = time.Now()
	body.CreatedBy = "test"
	body.Short = domain + shortLink
	body.ShortId = shortLink
	body.Clicks = 0
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
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, link)
}

func GetAllLinks(c *gin.Context) {
	links, err := repository.GetAllLinks()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, links)
}

func GenerateShortLink() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6

	shortLink := make([]byte, length)
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range shortLink {
		shortLink[i] = charset[rand.Intn(len(charset))]
	}

	return string(shortLink)
}
