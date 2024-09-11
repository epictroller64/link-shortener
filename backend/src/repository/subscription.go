package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type SubscriptionStatus string

const (
	SubscriptionStatusIncomplete        SubscriptionStatus = "incomplete"
	SubscriptionStatusIncompleteExpired SubscriptionStatus = "incomplete_expired"
	SubscriptionStatusTrialing          SubscriptionStatus = "trialing"
	SubscriptionStatusActive            SubscriptionStatus = "active"
	SubscriptionStatusPastDue           SubscriptionStatus = "past_due"
	SubscriptionStatusCanceled          SubscriptionStatus = "canceled"
	SubscriptionStatusUnpaid            SubscriptionStatus = "unpaid"
	SubscriptionStatusPaused            SubscriptionStatus = "paused"
)

type Subscription struct {
	ID               string             `json:"id"`
	CustomerID       string             `json:"customerId"`
	Status           SubscriptionStatus `json:"status"`
	CurrentPeriodEnd time.Time          `json:"currentPeriodEnd"`
	CreatedAt        time.Time          `json:"createdAt"`
	UpdatedAt        time.Time          `json:"updatedAt"`
	PackageID        string             `json:"packageId"`
}

// Package is a struct that represents a package that a user can subscribe to, no more then 3 needed
type Package struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Price         int    `json:"price"`
	MaxLinks      int    `json:"maxLinks"`
	MaxClicks     int    `json:"maxClicks"`
	CustomDomains int    `json:"customDomains"` // Will add this later if ever
	IsDefault     bool   `json:"isDefault"`
}

type Billing struct {
	Package      *Package      `json:"package"`
	Subscription *Subscription `json:"subscription"`
}

func GetPackageByID(id string) (Package, error) {
	query := `
		SELECT * FROM packages WHERE id = $1
	`

	row := Db.QueryRow(context.Background(), query, id)

	var subPackage Package
	err := row.Scan(&subPackage.ID, &subPackage.Name, &subPackage.Description, &subPackage.Price, &subPackage.MaxLinks, &subPackage.MaxClicks, &subPackage.CustomDomains, &subPackage.IsDefault)
	if err == pgx.ErrNoRows {
		return Package{}, nil
	}
	if err != nil {
		return Package{}, err
	}
	return subPackage, nil
}

func GetPackages() ([]Package, error) {
	query := `
		SELECT * FROM packages
	`

	rows, err := Db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var packages []Package
	for rows.Next() {
		var subPackage Package
		err := rows.Scan(&subPackage.ID, &subPackage.Name, &subPackage.Description, &subPackage.Price, &subPackage.MaxLinks, &subPackage.MaxClicks, &subPackage.CustomDomains, &subPackage.IsDefault)
		if err != nil {
			return nil, err
		}
		packages = append(packages, subPackage)
	}

	return packages, nil
}

func CreateSubscription(subscription Subscription) (Subscription, error) {
	query := `
		INSERT INTO subscriptions (id, customer_id, status, current_period_end, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := Db.Exec(context.Background(), query, subscription.ID, subscription.CustomerID, subscription.Status,
		subscription.CurrentPeriodEnd, subscription.CreatedAt, subscription.UpdatedAt)

	return subscription, err
}

// UpdateSubscription updates a subscription in the database, only the status, current_period_end and updated_at can be updated
func UpdateSubscription(subscription Subscription) (Subscription, error) {
	query := `
		UPDATE subscriptions
		SET status = $2, current_period_end = $3, updated_at = $4
		WHERE id = $1
	`

	_, err := Db.Exec(context.Background(), query, subscription.ID, subscription.Status, subscription.CurrentPeriodEnd, subscription.UpdatedAt)

	return subscription, err
}

// GetSubscriptionByCustomerId gets a subscription by customer id. This is assigned by stripe
func GetSubscriptionByCustomerId(customerId string) (*Subscription, error) {
	query := `
		SELECT * FROM subscriptions WHERE customer_id = $1
	`

	row := Db.QueryRow(context.Background(), query, customerId)

	var subscription Subscription
	err := row.Scan(&subscription.ID, &subscription.CustomerID, &subscription.Status, &subscription.CurrentPeriodEnd, &subscription.CreatedAt, &subscription.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil // Return nil, nil when no subscription is found
	}
	if err != nil {
		return nil, err // Return nil and the error for other errors
	}
	return &subscription, nil
}

func GetSubscriptionByID(id string) (Subscription, error) {
	query := `
		SELECT * FROM subscriptions WHERE id = $1
	`

	row := Db.QueryRow(context.Background(), query, id)

	var subscription Subscription
	err := row.Scan(&subscription.ID, &subscription.CustomerID, &subscription.Status, &subscription.CurrentPeriodEnd, &subscription.CreatedAt, &subscription.UpdatedAt)

	return subscription, err
}

func GetPackageById(id string) (*Package, error) {
	query := `
		SELECT * FROM packages WHERE id = $1
	`

	row := Db.QueryRow(context.Background(), query, id)

	var subPackage Package
	err := row.Scan(&subPackage.ID, &subPackage.Name, &subPackage.Description, &subPackage.Price, &subPackage.MaxLinks, &subPackage.MaxClicks, &subPackage.CustomDomains, &subPackage.IsDefault)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &subPackage, nil
}
