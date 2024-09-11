package handlers

import (
	"link-shortener-backend/src/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBilling(c *gin.Context) {
	user := c.MustGet("user").(*repository.User)

	// Function to return empty billing
	emptyBilling := func() {
		c.JSON(http.StatusOK, repository.Billing{
			Package:      nil,
			Subscription: nil,
		})
	}

	if user.StripeCustomerID == nil {
		emptyBilling()
		return
	}

	subscription, err := repository.GetSubscriptionByCustomerId(*user.StripeCustomerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if subscription == nil {
		emptyBilling()
		return
	}

	subPackage, err := repository.GetPackageById(subscription.PackageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, repository.Billing{
		Package:      subPackage,
		Subscription: subscription,
	})
}

func GetAccountDetails(c *gin.Context) {
	user := c.MustGet("user").(*repository.User)
	c.JSON(http.StatusOK, user)
}
