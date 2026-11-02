package analytics

import (
	"time"
)

// Click represents the click log stored in the database
type Click struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ShortCode  string    `gorm:"type:varchar(50);index;not null" json:"short_code"`
	IP         string    `gorm:"type:varchar(45);index" json:"ip"`
	Country    string    `gorm:"type:varchar(50)" json:"country"`
	City       string    `gorm:"type:varchar(100)" json:"city"`
	UserAgent  string    `gorm:"type:text" json:"user_agent"`
	DeviceType string    `gorm:"type:varchar(30);index" json:"device_type"`
	Browser    string    `gorm:"type:varchar(50);index" json:"browser"`
	OS         string    `gorm:"type:varchar(50);index" json:"os"`
	Referrer   string    `gorm:"type:text" json:"referrer"`
	ClickedAt  time.Time `gorm:"index;not null" json:"clicked_at"`
}

func (Click) TableName() string {
	return "clicks"
}

// Analytics Aggregate Models
type StatItem struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type TimeSeriesItem struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type AnalyticsSummary struct {
	ShortCode      string           `json:"short_code"`
	TotalClicks    int64            `json:"total_clicks"`
	UniqueVisitors int64            `json:"unique_visitors"`
	Devices        []StatItem       `json:"devices"`
	Browsers       []StatItem       `json:"browsers"`
	OperatingSys   []StatItem       `json:"operating_systems"`
	TopReferrers   []StatItem       `json:"top_referrers"`
	TopCountries   []StatItem       `json:"top_countries"`
	RecentClicks   []TimeSeriesItem `json:"recent_clicks"`
}
