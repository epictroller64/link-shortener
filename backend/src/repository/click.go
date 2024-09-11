package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mssola/useragent"
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

// GetClicks gets all clicks for a link
func GetClicks(linkId string) ([]Click, error) {
	query := `
		SELECT id, link_id, created_at, user_agent, referer, ip, country
		FROM clicks
		WHERE link_id = $1
	`

	var clicks []Click
	rows, err := Db.Query(context.Background(), query, linkId)
	if err != nil {
		return []Click{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var click Click
		err := rows.Scan(
			&click.ID,
			&click.LinkID,
			&click.CreatedAt,
			&click.UserAgent,
			&click.Referer,
			&click.IP,
			&click.Country,
		)
		if err != nil {
			return []Click{}, err
		}
		clicks = append(clicks, click)
	}

	return clicks, nil
}

type DailyStatistics struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

type DeviceType string

const (
	DeviceTypeDesktop DeviceType = "desktop"
	DeviceTypeMobile  DeviceType = "mobile"
)

type DeviceStatistics struct {
	Device DeviceType `json:"device"`
	Count  int        `json:"count"`
}

type TotalStatsResponse struct {
	TotalLinks  int `json:"totalLinks"`
	TotalClicks int `json:"totalClicks"`
}

func GetTotalStats(userId string) (*TotalStatsResponse, error) {
	query := `
		SELECT COUNT(*) as total_links, SUM(clicks) as total_clicks
		FROM links
		WHERE created_by = $1
	`

	var totalLinkCount int
	var totalClickCount int
	err := Db.QueryRow(context.Background(), query, userId).Scan(&totalLinkCount, &totalClickCount)
	if err != nil {
		return nil, err
	}

	return &TotalStatsResponse{TotalLinks: totalLinkCount, TotalClicks: totalClickCount}, nil
}

func GetDailyStatistics(userId string, startDate time.Time, endDate time.Time) ([]DailyStatistics, error) {
	query := `
		SELECT DATE(clicks.created_at) as date, COUNT(*) as count
		FROM clicks
		INNER JOIN links ON links.id = clicks.link_id
		WHERE links.created_by = $1 AND clicks.created_at BETWEEN $2 AND $3
		GROUP BY DATE(clicks.created_at)
		ORDER BY date
	`

	rows, err := Db.Query(context.Background(), query, userId, startDate, endDate)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []DailyStatistics{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	var dailyStatistics []DailyStatistics
	for rows.Next() {
		var stat DailyStatistics
		err := rows.Scan(&stat.Date, &stat.Count)
		if err != nil {
			return nil, err
		}
		dailyStatistics = append(dailyStatistics, stat)
	}

	return dailyStatistics, nil
}

func GetClicksByDateRange(linkId string, startDate time.Time, endDate time.Time) ([]Click, error) {
	query := `
		SELECT id, link_id, created_at, user_agent, referer, ip, country
		FROM clicks
		WHERE link_id = $1 AND created_at BETWEEN $2 AND $3
	`

	var clicks []Click = make([]Click, 0)
	rows, err := Db.Query(context.Background(), query, linkId, startDate, endDate)
	if err != nil {
		return clicks, err
	}
	defer rows.Close()

	for rows.Next() {
		var click Click
		err := rows.Scan(
			&click.ID,
			&click.LinkID,
			&click.CreatedAt,
			&click.UserAgent,
			&click.Referer,
			&click.IP,
			&click.Country,
		)
		if err != nil {
			return clicks, err
		}
		clicks = append(clicks, click)
	}

	return clicks, nil

}

func GetDeviceStatistics(userId string, linkId string, startDate time.Time, endDate time.Time) ([]DeviceStatistics, error) {
	query := `SELECT user_agent FROM clicks WHERE link_id = $1 AND created_at BETWEEN $2 AND $3`

	rows, err := Db.Query(context.Background(), query, linkId, startDate, endDate)
	var deviceStatistics []DeviceStatistics = make([]DeviceStatistics, 0)
	if err != nil {
		if err == pgx.ErrNoRows {
			return deviceStatistics, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userAgent string
		err := rows.Scan(&userAgent)
		if err != nil {
			return nil, err
		}
		device := parseUserAgent(userAgent)
		deviceStatistics = append(deviceStatistics, DeviceStatistics{Device: device, Count: 1})
	}

	return deviceStatistics, nil
}

func parseUserAgent(userAgentString string) DeviceType {
	ua := useragent.New(userAgentString)

	if ua.Mobile() {
		return DeviceTypeMobile
	}

	return DeviceTypeDesktop
}
