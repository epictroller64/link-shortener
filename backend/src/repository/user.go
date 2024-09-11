package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID               string    `json:"id"`
	Email            string    `json:"email"`
	Password         string    `json:"password"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	IpAddress        string    `json:"ipAddress"`
	UserAgent        string    `json:"userAgent"`
	StripeCustomerID *string   `json:"stripeCustomerID"`
}

func CreateUser(user User) error {
	query := `
		INSERT INTO users (email, password, created_at, updated_at, ip_address, user_agent, stripe_customer_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := Db.Exec(context.Background(), query, user.Email, user.Password, user.CreatedAt, user.UpdatedAt, user.IpAddress, user.UserAgent, user.StripeCustomerID)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(email string) (*User, error) {
	query := `
		SELECT * FROM users WHERE email = $1
	`
	fmt.Println(email)
	row := Db.QueryRow(context.Background(), query, email)
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.IpAddress, &user.UserAgent, &user.StripeCustomerID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		fmt.Println(err.Error())
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id string) (*User, error) {
	query := `
		SELECT * FROM users WHERE id = $1
	`
	row := Db.QueryRow(context.Background(), query, id)
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.IpAddress, &user.UserAgent, &user.StripeCustomerID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
