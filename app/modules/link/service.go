package link

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	appErrors "link-shortener/common/errors"
	"link-shortener/common/utils"
	"link-shortener/config"
)

type Service interface {
	CreateLink(ctx context.Context, userID uint, req CreateLinkRequest) (*LinkResponse, error)
	GetLinkByID(ctx context.Context, id, userID uint) (*LinkResponse, error)
	ListLinks(ctx context.Context, userID uint, page, limit int) ([]LinkResponse, int64, error)
	UpdateLink(ctx context.Context, id, userID uint, req UpdateLinkRequest) (*LinkResponse, error)
	DeleteLink(ctx context.Context, id, userID uint) error
}

type service struct {
	repo  Repository
	cache RedisCache
	cfg   *config.AppConfig
}

func NewService(repo Repository, cache RedisCache, cfg *config.AppConfig) Service {
	return &service{
		repo:  repo,
		cache: cache,
		cfg:   cfg,
	}
}

func (s *service) CreateLink(ctx context.Context, userID uint, req CreateLinkRequest) (*LinkResponse, error) {
	parsedURL, err := url.ParseRequestURI(strings.TrimSpace(req.OriginalURL))
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return nil, appErrors.BadRequest("invalid URL format, must begin with http:// or https://", nil)
	}

	var shortCode string
	if req.CustomCode != "" {
		shortCode = strings.TrimSpace(req.CustomCode)
		existing, err := s.repo.FindByShortCode(ctx, shortCode)
		if err != nil {
			return nil, appErrors.Internal("failed to check short code availability")
		}
		if existing != nil {
			return nil, appErrors.Conflict("custom short code is already taken")
		}
	} else {
		// Generate random 6-character short code
		for i := 0; i < 5; i++ {
			code, err := utils.GenerateShortCode(6)
			if err != nil {
				return nil, appErrors.Internal("failed to generate short code")
			}
			existing, err := s.repo.FindByShortCode(ctx, code)
			if err != nil {
				return nil, appErrors.Internal("failed to check code uniqueness")
			}
			if existing == nil {
				shortCode = code
				break
			}
		}
		if shortCode == "" {
			return nil, appErrors.Internal("could not generate unique short code, please try again")
		}
	}

	link := &Link{
		UserID:      userID,
		OriginalURL: parsedURL.String(),
		ShortCode:   shortCode,
		Title:       req.Title,
		Description: req.Description,
		IsActive:    true,
		ExpiresAt:   req.ExpiresAt,
	}

	if err := s.repo.Create(ctx, link); err != nil {
		return nil, appErrors.Internal("failed to save link")
	}

	// Warm cache immediately in Redis for ultra-fast first redirect
	_ = s.cache.SetLink(ctx, shortCode, link.OriginalURL, 24*time.Hour)
	_ = s.cache.SetLinkMeta(ctx, shortCode, link, 24*time.Hour)

	return s.toResponse(link), nil
}

func (s *service) GetLinkByID(ctx context.Context, id, userID uint) (*LinkResponse, error) {
	link, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, appErrors.Internal("failed to query link")
	}
	if link == nil || link.UserID != userID {
		return nil, appErrors.NotFound("link not found")
	}
	return s.toResponse(link), nil
}

func (s *service) ListLinks(ctx context.Context, userID uint, page, limit int) ([]LinkResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	links, total, err := s.repo.ListByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, appErrors.Internal("failed to list links")
	}

	responses := make([]LinkResponse, len(links))
	for i, l := range links {
		responses[i] = *s.toResponse(&l)
	}

	return responses, total, nil
}

func (s *service) UpdateLink(ctx context.Context, id, userID uint, req UpdateLinkRequest) (*LinkResponse, error) {
	link, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, appErrors.Internal("failed to query link")
	}
	if link == nil || link.UserID != userID {
		return nil, appErrors.NotFound("link not found")
	}

	if req.OriginalURL != "" {
		parsedURL, err := url.ParseRequestURI(strings.TrimSpace(req.OriginalURL))
		if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
			return nil, appErrors.BadRequest("invalid URL format", nil)
		}
		link.OriginalURL = parsedURL.String()
	}
	if req.Title != "" {
		link.Title = req.Title
	}
	if req.Description != "" {
		link.Description = req.Description
	}
	if req.IsActive != nil {
		link.IsActive = *req.IsActive
	}
	if req.ExpiresAt != nil {
		link.ExpiresAt = req.ExpiresAt
	}

	if err := s.repo.Update(ctx, link); err != nil {
		return nil, appErrors.Internal("failed to update link")
	}

	// Invalidate or refresh Redis cache
	if !link.IsActive || (link.ExpiresAt != nil && link.ExpiresAt.Before(time.Now())) {
		_ = s.cache.DeleteLink(ctx, link.ShortCode)
	} else {
		_ = s.cache.SetLink(ctx, link.ShortCode, link.OriginalURL, 24*time.Hour)
		_ = s.cache.SetLinkMeta(ctx, link.ShortCode, link, 24*time.Hour)
	}

	return s.toResponse(link), nil
}

func (s *service) DeleteLink(ctx context.Context, id, userID uint) error {
	link, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return appErrors.Internal("failed to query link")
	}
	if link == nil || link.UserID != userID {
		return appErrors.NotFound("link not found")
	}

	if err := s.repo.Delete(ctx, id, userID); err != nil {
		return appErrors.Internal("failed to delete link")
	}

	_ = s.cache.DeleteLink(ctx, link.ShortCode)
	return nil
}

func (s *service) toResponse(link *Link) *LinkResponse {
	return &LinkResponse{
		ID:          link.ID,
		OriginalURL: link.OriginalURL,
		ShortCode:   link.ShortCode,
		ShortURL:    fmt.Sprintf("%s/%s", strings.TrimRight(s.cfg.App.BaseURL, "/"), link.ShortCode),
		Title:       link.Title,
		Description: link.Description,
		IsActive:    link.IsActive,
		ExpiresAt:   link.ExpiresAt,
		ClickCount:  link.ClickCount,
		CreatedAt:   link.CreatedAt,
		UpdatedAt:   link.UpdatedAt,
	}
}
