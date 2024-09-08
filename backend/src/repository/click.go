package repository

import (
	"context"
	"time"
)

// TrackClick tracks a click on a link
type Click struct {
	ID        int       `json:"id,omitempty"`
	LinkID    int       `json:"linkId"`
	CreatedAt time.Time `json:"createdAt"`
	UserAgent string    `json:"userAgent"`
	Referer   string    `json:"referer"`
	IP        string    `json:"ip"`
	Country   string    `json:"country"`
}

func CreateClick(click Click) (Click, error) {
	query := `
		INSERT INTO clicks (link_id, created_at, user_agent, referer, ip, country)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, link_id, created_at, user_agent, referer, ip, country
	`

	err := Db.QueryRow(
		context.Background(),
		query,
		click.LinkID,
		click.CreatedAt,
		click.UserAgent,
		click.Referer,
		click.IP,
		click.Country,
	).Scan(
		&click.ID,
		&click.LinkID,
		&click.CreatedAt,
		&click.UserAgent,
		&click.Referer,
		&click.IP,
		&click.Country,
	)

	if err != nil {
		return Click{}, err
	}

	return click, nil
}
