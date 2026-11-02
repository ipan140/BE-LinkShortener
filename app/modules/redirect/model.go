package redirect

import (
	"time"
)

// ClickData contains all telemetry collected on a link click
type ClickData struct {
	ShortCode  string    `json:"short_code"`
	IP         string    `json:"ip"`
	Country    string    `json:"country"`
	City       string    `json:"city"`
	UserAgent  string    `json:"user_agent"`
	DeviceType string    `json:"device_type"`
	Browser    string    `json:"browser"`
	OS         string    `json:"os"`
	Referrer   string    `json:"referrer"`
	ClickedAt  time.Time `json:"clicked_at"`
}

// RedirectLink is a lightweight model representation for redirect evaluation
type RedirectLink struct {
	ID          uint
	OriginalURL string
	ShortCode   string
	IsActive    bool
	ExpiresAt   *time.Time
}
