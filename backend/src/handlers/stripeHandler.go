package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"link-shortener-backend/src/repository"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	stripe "github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
	"github.com/stripe/stripe-go/v79/subscription"
	"github.com/stripe/stripe-go/webhook"
)

func init() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}

var (
	ErrSubscriptionExists = errors.New("subscription already exists")
)

type StripeCheckoutSession struct {
	Price    string `json:"price"`
	Quantity int64  `json:"quantity"`
}

func StripeCreateCheckoutSession(c *gin.Context) {

	var checkoutSessionBody StripeCheckoutSession
	err := c.BindJSON(&checkoutSessionBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(checkoutSessionBody.Price),
				Quantity: stripe.Int64(checkoutSessionBody.Quantity),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String("http://localhost:3000/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String("http://localhost:3000/cancel"),
	}

	checkoutSession, err := session.New(params)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusSeeOther, checkoutSession.URL)
}
func StripeSuccess(c *gin.Context) {
	fmt.Println("Success")
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func StripeWebHook(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	signature := c.Request.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, signature, os.Getenv("STRIPE_WEBHOOK_SECRET"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		CreatePayment(&paymentIntent)
	case "customer.subscription.created":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//CreateSubscription(&subscription)
	case "checkout.session.completed":
		// Provision the subscription here
		var checkoutSession stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &checkoutSession)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("checkout session completed")
		customerID := checkoutSession.Customer.ID
		subscriptionID := checkoutSession.Subscription.ID
		subscription, err := subscription.Get(subscriptionID, &stripe.SubscriptionParams{
			Customer: &customerID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = CreateSubscription(subscription)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//CreateCheckoutSession(&checkoutSession)
	case "customer.subscription.updated":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		UpdateSubscription(&subscription)
	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//DeleteSubscription(&subscription)
	case "charge.failed":
		var charge stripe.Charge
		err := json.Unmarshal(event.Data.Raw, &charge)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "subscription_schedule.expiring":
		// When 7 days are left
		var subscriptionSchedule stripe.SubscriptionSchedule
		err := json.Unmarshal(event.Data.Raw, &subscriptionSchedule)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//UpdateSubscriptionSchedule(&subscriptionSchedule)
	case "invoice.payment_failed":
		// Invoice payment failed
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//UpdateInvoice(&invoice)
	default:
		c.JSON(http.StatusOK, gin.H{"message": "Webhook received"})
	}
}

func CreateSubscription(stripeSubscription *stripe.Subscription) error {
	subscription := &repository.Subscription{
		ID:               stripeSubscription.ID,
		CustomerID:       stripeSubscription.Customer.ID,
		Status:           repository.SubscriptionStatus(stripeSubscription.Status),
		CurrentPeriodEnd: time.Unix(stripeSubscription.CurrentPeriodEnd, 0),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	existingSubscription, err := repository.GetSubscriptionByCustomerId(subscription.CustomerID)
	if err != nil {
		return err
	}
	if existingSubscription != nil {
		return ErrSubscriptionExists
	}

	_, err = repository.CreateSubscription(*subscription)

	return err
}

func UpdateSubscription(stripeSubscription *stripe.Subscription) error {
	subscription := &repository.Subscription{
		ID:               stripeSubscription.ID,
		CustomerID:       stripeSubscription.Customer.ID,
		Status:           repository.SubscriptionStatus(stripeSubscription.Status),
		CurrentPeriodEnd: time.Unix(stripeSubscription.CurrentPeriodEnd, 0),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	_, err := repository.UpdateSubscription(*subscription)

	return err
}

func CreatePayment(stripePayment *stripe.PaymentIntent) error {
	payment := &repository.Payment{
		ID:        stripePayment.ID,
		Amount:    stripePayment.Amount,
		Currency:  string(stripePayment.Currency),
		Status:    string(stripePayment.Status),
		CreatedAt: time.Unix(stripePayment.Created, 0),
		UpdatedAt: time.Unix(stripePayment.Created, 0),
	}

	_, err := repository.CreatePayment(*payment)

	return err
}

// Cron job to sync subscriptions with stripe
func StripeSubscriptionSync(c *gin.Context) {
	iter := subscription.List(&stripe.SubscriptionListParams{
		Status: stripe.String(string(stripe.SubscriptionStatusActive)),
	})
	for iter.Next() {
		subscription := iter.Subscription()
		subscriptionJSON, err := json.MarshalIndent(subscription, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling subscription:", err)
			continue
		}
		subscriptionItems := subscription.Items
		for range subscriptionItems.Data {
			dbSubscription, err := repository.GetSubscriptionByID(subscription.ID)
			if err != nil {
				fmt.Println("Error getting subscription by ID:", err)
				continue
			}
			if dbSubscription.Status != repository.SubscriptionStatus(subscription.Status) {
				fmt.Println("Subscription status changed:", dbSubscription.Status, "->", repository.SubscriptionStatus(subscription.Status))
				// Send email to user
				err = UpdateSubscription(subscription)
				if err != nil {
					fmt.Println("Error updating subscription:", err)
				}
			}

		}
		fmt.Println(string(subscriptionJSON))
	}
}
