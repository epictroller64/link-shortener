package repository

import (
	"context"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IpAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
}

func CreateUser(user User) error {
	query := `
		INSERT INTO users (email, password, created_at, updated_at, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := Db.Exec(context.Background(), query, user.Email, user.Password, user.CreatedAt, user.UpdatedAt, user.IpAddress, user.UserAgent)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(email string) (*User, error) {
	query := `
		SELECT * FROM users WHERE email = $1
	`
	row := Db.QueryRow(context.Background(), query, email)
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.IpAddress, &user.UserAgent)
	if err != nil {
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
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.IpAddress, &user.UserAgent)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
