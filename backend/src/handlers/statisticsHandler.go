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

type TotalStatsResponse struct {
	TotalLinks  int `json:"totalLinks"`
	TotalClicks int `json:"totalClicks"`
}

func GetTotalStats(c *gin.Context) {
	user := c.MustGet("user").(*repository.User)
	totalStats, err := repository.GetTotalStats(user.ID)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, totalStats)
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
	stats, err := repository.GetDailyStatistics(user.ID, startDate, endDate)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(stats)
	c.JSON(http.StatusOK, stats)
}

// GetStatistics returns the number of clicks for each day in the given date range. This is targeting specific link
func GetStatistics(c *gin.Context) {
	user := c.MustGet("user").(*repository.User)
	// Check whether user owns this link at all
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

	link, err := repository.GetLink(request.LinkId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if link == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}
	if link.CreatedBy != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this link"})
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

func GetDeviceStatistics(c *gin.Context) {
	user := c.MustGet("user").(*repository.User)
	var request StatisticsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	start, end, err := ParseDates(request.StartDate, request.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	stats, err := repository.GetDeviceStatistics(user.ID, request.LinkId, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func ParseDates(startDate string, endDate string) (time.Time, time.Time, error) {
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return start, end, nil
}

func GetRefererStatistics(c *gin.Context) {
	user := c.MustGet("user").(*repository.User)
	var request StatisticsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(request)
	start, end, err := ParseDates(request.StartDate, request.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	isOwned, err := CheckLinkOwnership(request.LinkId, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isOwned {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this link"})
		return
	}
	stats, err := repository.GetRefererStatistics(request.LinkId, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func GetIpStatistics(c *gin.Context) {
	user := c.MustGet("user").(*repository.User)
	var request StatisticsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	isOwned, err := CheckLinkOwnership(request.LinkId, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isOwned {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this link"})
		return
	}
	start, end, err := ParseDates(request.StartDate, request.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	stats, err := repository.GetIpStatistics(request.LinkId, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func CheckLinkOwnership(linkId string, user *repository.User) (bool, error) {
	link, err := repository.GetLink(linkId)
	if err != nil {
		return false, err
	}
	if link == nil {
		return false, nil
	}
	if link.CreatedBy != user.ID {
		return false, nil
	}
	return true, nil
}
