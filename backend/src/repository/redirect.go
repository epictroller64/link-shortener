package repository

import "context"

type TargetType string

const (
	Header TargetType = "header"
	Cookie TargetType = "cookie"
)

type TargetMethod string

const (
	Match      TargetMethod = "match"
	Regex      TargetMethod = "regex"
	Contains   TargetMethod = "contains"
	StartsWith TargetMethod = "startsWith"
	EndsWith   TargetMethod = "endsWith"
)

type Redirect struct {
	ID           int          `json:"id"`
	LinkID       int          `json:"linkID"`
	TargetType   TargetType   `json:"targetType"`
	TargetMethod TargetMethod `json:"targetMethod"`
	RedirectURL  string       `json:"redirectURL"`
	TargetValue  *string      `json:"targetValue"`
	TargetName   *string      `json:"targetName"`
}

func CreateRedirect(redirect Redirect) (Redirect, error) {
	query := `
		INSERT INTO redirects (link_id, target_type, target_method, redirect_url, target_value, target_name)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *
	`

	row := Db.QueryRow(context.Background(), query, redirect.LinkID, redirect.TargetType, redirect.TargetMethod, redirect.RedirectURL, redirect.TargetValue, redirect.TargetName)

	var createdRedirect Redirect
	err := row.Scan(&createdRedirect.ID, &createdRedirect.LinkID, &createdRedirect.TargetType, &createdRedirect.TargetMethod, &createdRedirect.RedirectURL, &createdRedirect.TargetValue, &createdRedirect.TargetName)
	if err != nil {
		return Redirect{}, err
	}

	return createdRedirect, nil
}

func GetRedirectsByLinkID(linkID string) ([]Redirect, error) {
	query := `
		SELECT * FROM redirects WHERE link_id = $1
	`

	rows, err := Db.Query(context.Background(), query, linkID)
	var redirects []Redirect = make([]Redirect, 0)
	if err != nil {
		return redirects, err
	}

	for rows.Next() {
		var redirect Redirect
		err := rows.Scan(&redirect.ID, &redirect.LinkID, &redirect.TargetType, &redirect.TargetMethod, &redirect.RedirectURL, &redirect.TargetValue, &redirect.TargetName)
		if err != nil {
			return redirects, err
		}
		redirects = append(redirects, redirect)
	}

	return redirects, nil
}

func DeleteRedirect(redirectID string) error {
	query := `
		DELETE FROM redirects WHERE id = $1
	`

	_, err := Db.Exec(context.Background(), query, redirectID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateRedirect(redirect Redirect) (Redirect, error) {
	query := `
		UPDATE redirects SET target_type = $1, target_method = $2, redirect_url = $3, target_value = $4, target_name = $5 WHERE id = $6
		RETURNING *
	`

	row := Db.QueryRow(context.Background(), query, redirect.TargetType, redirect.TargetMethod, redirect.RedirectURL, redirect.TargetValue, redirect.TargetName, redirect.ID)

	var updatedRedirect Redirect
	err := row.Scan(&updatedRedirect.ID, &updatedRedirect.LinkID, &updatedRedirect.TargetType, &updatedRedirect.TargetMethod, &updatedRedirect.RedirectURL)
	if err != nil {
		return Redirect{}, err
	}

	return updatedRedirect, nil
}
