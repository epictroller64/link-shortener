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
}

func CreateClick(click Click) (Click, error) {
	query := `
		INSERT INTO clicks (id, link_id, created_at)
		VALUES ($1, $2, $3)
		RETURNING id, link_id, created_at
	`

	err := Db.QueryRow(
		context.Background(),
		query,
		click.ID,
		click.LinkID,
		click.CreatedAt,
	).Scan(
		&click.ID,
		&click.LinkID,
		&click.CreatedAt,
	)

	if err != nil {
		return Click{}, err
	}

	return click, nil
}
