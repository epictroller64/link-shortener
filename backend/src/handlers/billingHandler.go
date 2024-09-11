package handlers

import (
	"link-shortener-backend/src/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBilling(c *gin.Context) {
	user := c.MustGet("user").(*repository.User)
	subscription, err := repository.GetSubscriptionByCustomerId(user.StripeCustomerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	subPackage, err := repository.GetPackageById(subscription.PackageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	billing := repository.Billing{
		Package:      *subPackage,
		Subscription: *subscription,
	}
	c.JSON(http.StatusOK, billing)
}

func GetAccountDetails(c *gin.Context) {
	user := c.MustGet("user").(*repository.User)
	c.JSON(http.StatusOK, user)
}
