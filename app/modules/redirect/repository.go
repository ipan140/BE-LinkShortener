package redirect

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"link-shortener/app/modules/link"
)

type Repository interface {
	FindActiveLink(ctx context.Context, shortCode string) (*RedirectLink, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindActiveLink(ctx context.Context, shortCode string) (*RedirectLink, error) {
	var l link.Link
	// Select only necessary columns to minimize memory allocation and transfer
	err := r.db.WithContext(ctx).
		Select("id, original_url, short_code, is_active, expires_at").
		Where("short_code = ?", shortCode).
		First(&l).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &RedirectLink{
		ID:          l.ID,
		OriginalURL: l.OriginalURL,
		ShortCode:   l.ShortCode,
		IsActive:    l.IsActive,
		ExpiresAt:   l.ExpiresAt,
	}, nil
}
