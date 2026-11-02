package link

import (
	"time"

	"gorm.io/gorm"
)

// Link represents a shortened URL entity in the database
type Link struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	OriginalURL  string         `gorm:"type:text;not null" json:"original_url"`
	IOSURL       string         `gorm:"type:text" json:"ios_url,omitempty"`
	AndroidURL   string         `gorm:"type:text" json:"android_url,omitempty"`
	ShortCode    string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"short_code"`
	Title        string         `gorm:"type:varchar(255)" json:"title"`
	Description  string         `gorm:"type:text" json:"description"`
	ImageURL     string         `gorm:"type:text" json:"image_url,omitempty"` // OpenGraph preview image banner
	Tags         string         `gorm:"type:varchar(255);index" json:"tags,omitempty"` // Campaign tags (e.g. marketing,promo)
	IsActive     bool           `gorm:"default:true;index" json:"is_active"`
	ExpiresAt    *time.Time     `gorm:"index" json:"expires_at,omitempty"`
	MaxClicks    *int64         `json:"max_clicks,omitempty"`
	ClickCount   int64          `gorm:"default:0" json:"click_count"`
	PasswordHash string         `gorm:"type:varchar(255)" json:"-"`
	IsProtected  bool           `gorm:"default:false" json:"is_protected"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Link) TableName() string {
	return "links"
}

// Request and Response DTOs
type CreateLinkRequest struct {
	OriginalURL string     `json:"original_url" validate:"required,url"`
	IOSURL      string     `json:"ios_url,omitempty" validate:"omitempty,url"`
	AndroidURL  string     `json:"android_url,omitempty" validate:"omitempty,url"`
	CustomCode  string     `json:"custom_code,omitempty" validate:"omitempty,alphanum,min=3,max=30"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	ImageURL    string     `json:"image_url,omitempty"`
	Tags        string     `json:"tags,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	MaxClicks   *int64     `json:"max_clicks,omitempty"`
	Password    string     `json:"password,omitempty"`
}

type UpdateLinkRequest struct {
	OriginalURL string     `json:"original_url,omitempty"`
	IOSURL      *string    `json:"ios_url,omitempty"`
	AndroidURL  *string    `json:"android_url,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	ImageURL    *string    `json:"image_url,omitempty"`
	Tags        *string    `json:"tags,omitempty"`
	IsActive    *bool      `json:"is_active,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	MaxClicks   *int64     `json:"max_clicks,omitempty"`
	Password    *string    `json:"password,omitempty"` // empty string removes protection
}

type LinkResponse struct {
	ID          uint       `json:"id"`
	OriginalURL string     `json:"original_url"`
	IOSURL      string     `json:"ios_url,omitempty"`
	AndroidURL  string     `json:"android_url,omitempty"`
	ShortCode   string     `json:"short_code"`
	ShortURL    string     `json:"short_url"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ImageURL    string     `json:"image_url,omitempty"`
	Tags        string     `json:"tags,omitempty"`
	IsActive    bool       `json:"is_active"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	MaxClicks   *int64     `json:"max_clicks,omitempty"`
	ClickCount  int64      `json:"click_count"`
	IsProtected bool       `json:"is_protected"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
