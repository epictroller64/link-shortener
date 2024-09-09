package handlers

import (
	"fmt"
	"link-shortener-backend/src/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type StatisticsRequest struct {
	LinkId    string `json:"linkId"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type DailyStatisticsResponse struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// GetDailyStatistics returns the number of clicks for each day in the given date range. This is targeting whole account
func GetDailyStatistics(c *gin.Context) {
	user := c.MustGet("user").(*repository.User)
	var request DailyStatisticsResponse
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Parse the date strings into time.Time
	startDate, err := time.Parse(time.RFC3339, request.StartDate)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, request.EndDate)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}
	stats, err := repository.GetDailyStatistics(user.ID, startDate, endDate)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
func GetStatistics(c *gin.Context) {
	var request StatisticsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Parse the date strings into time.Time
	startDate, err := time.Parse(time.RFC3339, request.StartDate)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, request.EndDate)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}
	fmt.Println(startDate, endDate)
	clicks, err := repository.GetClicksByDateRange(request.LinkId, startDate, endDate)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, clicks)
}
