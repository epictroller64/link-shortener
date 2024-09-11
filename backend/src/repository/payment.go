package repository

import (
	"context"
	"time"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
)

type Payment struct {
	ID        string    `json:"id"`
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreatePayment(payment Payment) (Payment, error) {
	query := `
		INSERT INTO payments (id, amount, currency, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := Db.Exec(context.Background(), query, payment.ID, payment.Amount, payment.Currency, payment.Status, payment.CreatedAt, payment.UpdatedAt)

	return payment, err
}
