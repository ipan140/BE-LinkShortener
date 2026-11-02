package link

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, link *Link) error
	FindByID(ctx context.Context, id uint) (*Link, error)
	FindByShortCode(ctx context.Context, code string) (*Link, error)
	ListByUserID(ctx context.Context, userID uint, page, limit int) ([]Link, int64, error)
	Update(ctx context.Context, link *Link) error
	Delete(ctx context.Context, id, userID uint) error
	IncrementClickCount(ctx context.Context, linkID uint, delta int64) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, link *Link) error {
	return r.db.WithContext(ctx).Create(link).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*Link, error) {
	var link Link
	err := r.db.WithContext(ctx).First(&link, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &link, nil
}

func (r *repository) FindByShortCode(ctx context.Context, code string) (*Link, error) {
	var link Link
	err := r.db.WithContext(ctx).Where("short_code = ?", code).First(&link).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &link, nil
}

func (r *repository) ListByUserID(ctx context.Context, userID uint, page, limit int) ([]Link, int64, error) {
	var links []Link
	var total int64

	db := r.db.WithContext(ctx).Model(&Link{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&links).Error
	if err != nil {
		return nil, 0, err
	}

	return links, total, nil
}

func (r *repository) Update(ctx context.Context, link *Link) error {
	return r.db.WithContext(ctx).Save(link).Error
}

func (r *repository) Delete(ctx context.Context, id, userID uint) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&Link{}).Error
}

func (r *repository) IncrementClickCount(ctx context.Context, linkID uint, delta int64) error {
	return r.db.WithContext(ctx).Model(&Link{}).
		Where("id = ?", linkID).
		UpdateColumn("click_count", gorm.Expr("click_count + ?", delta)).Error
}
