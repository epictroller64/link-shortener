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
		SELECT COUNT(*) as total_links, COALESCE(SUM(clicks), 0) as total_clicks
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

// Daily statistics for the whole account, grouped by day
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
	var dailyStatistics []DailyStatistics = make([]DailyStatistics, 0)
	if err != nil {
		if err == pgx.ErrNoRows {
			return dailyStatistics, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stat DailyStatistics
		err := rows.Scan(&stat.Date, &stat.Count)
		if err != nil {
			return make([]DailyStatistics, 0), err
		}
		dailyStatistics = append(dailyStatistics, stat)
	}

	return dailyStatistics, nil
}

func GetClicksByDateRange(linkId string, startDate time.Time, endDate time.Time) ([]DailyStatistics, error) {
	query := `
		SELECT DATE(created_at) as date, COUNT(*) as count
		FROM clicks
		WHERE link_id = $1 AND created_at BETWEEN $2 AND $3
		GROUP BY DATE(created_at)
		ORDER BY date
	`

	rows, err := Db.Query(context.Background(), query, linkId, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dailyClicks []DailyStatistics
	for rows.Next() {
		var clickCount DailyStatistics
		err := rows.Scan(&clickCount.Date, &clickCount.Count)
		if err != nil {
			return nil, err
		}
		dailyClicks = append(dailyClicks, clickCount)
	}

	return dailyClicks, nil
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

type RefererStatistics struct {
	Referer string `json:"referer"`
	Count   int    `json:"count"`
}

func GetRefererStatistics(linkId string, startDate time.Time, endDate time.Time) ([]RefererStatistics, error) {
	query := `
		SELECT referer, COUNT(*) as count
		FROM clicks
		WHERE link_id = $1 AND created_at BETWEEN $2 AND $3
		GROUP BY referer
		ORDER BY count DESC
	`

	rows, err := Db.Query(context.Background(), query, linkId, startDate, endDate)
	if err != nil {
		if err == pgx.ErrNoRows {
			return make([]RefererStatistics, 0), nil
		}
		return nil, err
	}
	defer rows.Close()

	var refererStatistics []RefererStatistics = make([]RefererStatistics, 0)
	for rows.Next() {
		var referer string
		var count int
		err := rows.Scan(&referer, &count)
		if err != nil {
			return nil, err
		}
		refererStatistics = append(refererStatistics, RefererStatistics{Referer: referer, Count: count})
	}

	return refererStatistics, nil
}

type IpStatistics struct {
	Ip    string `json:"ip"`
	Count int    `json:"count"`
}

func GetIpStatistics(linkId string, startDate time.Time, endDate time.Time) ([]IpStatistics, error) {
	query := `
		SELECT ip, COUNT(*) as count
		FROM clicks
		WHERE link_id = $1 AND created_at BETWEEN $2 AND $3
		GROUP BY ip
		ORDER BY count DESC
	`

	rows, err := Db.Query(context.Background(), query, linkId, startDate, endDate)
	if err != nil {
		if err == pgx.ErrNoRows {
			return make([]IpStatistics, 0), nil
		}
		return nil, err
	}
	defer rows.Close()

	var ipStatistics []IpStatistics = make([]IpStatistics, 0)

	for rows.Next() {
		var ip string
		var count int
		err := rows.Scan(&ip, &count)
		if err != nil {
			return make([]IpStatistics, 0), err
		}
		ipStatistics = append(ipStatistics, IpStatistics{Ip: ip, Count: count})
	}

	return ipStatistics, nil

}
