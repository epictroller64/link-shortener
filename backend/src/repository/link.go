package repository

import (
	"context"
	"time"
)

type Link struct {
	ID        int       `json:"id,omitempty"`
	ShortId   string    `json:"shortId"`
	Original  string    `json:"original"`
	Short     string    `json:"short,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	CreatedBy string    `json:"createdBy,omitempty"`
	Clicks    int       `json:"clicks"`
}

// CreateLink creates a new link in the database
func CreateLink(link Link) (Link, error) {
	query := `
		INSERT INTO links (original, short, created_at, created_by, clicks, short_id)
		VALUES ($1, $2, $3, $4, 0, $5)
		RETURNING id, original, short, created_at, created_by, clicks, short_id
	`

	err := Db.QueryRow(
		context.Background(),
		query,
		link.Original,
		link.Short,
		link.CreatedAt,
		link.CreatedBy,
		link.ShortId,
	).Scan(
		&link.ID,
		&link.Original,
		&link.Short,
		&link.CreatedAt,
		&link.CreatedBy,
		&link.Clicks,
		&link.ShortId,
	)

	if err != nil {
		return Link{}, err
	}

	return link, nil
}

func GetLinkByShortId(shortId string) (Link, error) {
	query := `
		SELECT id, original, short, created_at, created_by, clicks, short_id
		FROM links
		WHERE short_id = $1
	`

	var link Link
	err := Db.QueryRow(context.Background(), query, shortId).Scan(
		&link.ID,
		&link.Original,
		&link.Short,
		&link.CreatedAt,
		&link.CreatedBy,
		&link.Clicks,
		&link.ShortId,
	)

	if err != nil {
		return Link{}, err
	}

	return link, nil
}

func GetLink(id string) (Link, error) {
	query := `
		SELECT id, original, short, created_at, created_by, clicks, short_id
		FROM links
		WHERE id = $1
	`

	var link Link
	err := Db.QueryRow(context.Background(), query, id).Scan(
		&link.ID,
		&link.Original,
		&link.Short,
		&link.CreatedAt,
		&link.CreatedBy,
		&link.Clicks,
		&link.ShortId,
	)

	if err != nil {
		return Link{}, err
	}

	return link, nil
}

func GetAllLinks() ([]Link, error) {
	query := `
		SELECT id, original, short, created_at, created_by, clicks, short_id
		FROM links ORDER BY id DESC
	`

	var links []Link = make([]Link, 0)
	rows, err := Db.Query(context.Background(), query)
	if err != nil {
		return links, err
	}
	defer rows.Close()

	for rows.Next() {
		var link Link
		err := rows.Scan(
			&link.ID,
			&link.Original,
			&link.Short,
			&link.CreatedAt,
			&link.CreatedBy,
			&link.Clicks,
			&link.ShortId,
		)
		if err != nil {
			return links, err
		}
		links = append(links, link)
	}

	return links, nil
}

// UpdateLinkClickCount increments the click count for a specific link
func UpdateLinkClickCount(linkID int) error {
	query := `
		UPDATE links
		SET clicks = clicks + 1
		WHERE id = $1
	`

	_, err := Db.Exec(context.Background(), query, linkID)
	if err != nil {
		return err
	}

	return nil
}
