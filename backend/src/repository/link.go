package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
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

func DeleteLink(id string, userID string) error {
	query := `
		DELETE FROM links
		WHERE id = $1 AND created_by = $2
	`

	_, err := Db.Exec(context.Background(), query, id, userID)
	if err != nil {
		return err
	}

	return nil
}

func GetLinkByShortId(shortId string) (*Link, error) {
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
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &link, nil
}

func GetLink(id string) (*Link, error) {
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
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &link, nil
}

func GetAllLinks(userID string) ([]Link, error) {
	query := `
		SELECT id, original, short, created_at, created_by, clicks, short_id
		FROM links
		WHERE created_by = $1
		ORDER BY id DESC
	`

	var links []Link = make([]Link, 0)
	rows, err := Db.Query(context.Background(), query, userID)

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

func GetRecentLinks(userID string) ([]Link, error) {
	query := `
		SELECT id, original, short, created_at, created_by, clicks, short_id
		FROM links
		WHERE created_by = $1
		ORDER BY created_at DESC
		LIMIT 10
	`

	var links []Link = make([]Link, 0)
	rows, err := Db.Query(context.Background(), query, userID)
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
