package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
	appErrors "link-shortener/common/errors"
	"link-shortener/common/middleware"
	"link-shortener/config"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	GetProfile(ctx context.Context, userID uint) (*User, error)
	RegenerateAPIKey(ctx context.Context, userID uint) (string, error)
	ValidateAPIKey(ctx context.Context, apiKey string) (*User, error)
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

func generateAPIKey() string {
	bytes := make([]byte, 16)
	_, _ = rand.Read(bytes)
	return fmt.Sprintf("ak_live_%s", hex.EncodeToString(bytes))
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || len(req.Password) < 6 {
		return nil, appErrors.BadRequest("valid email and password (min 6 characters) are required", nil)
	}

	existingUser, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, appErrors.Internal("failed to check existing user")
	}
	if existingUser != nil {
		return nil, appErrors.Conflict("email is already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, appErrors.Internal("failed to hash password")
	}

	apiKey := generateAPIKey()
	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user",
		APIKey:   apiKey,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, appErrors.Internal("failed to create user")
	}

	token, err := middleware.GenerateToken(user.ID, user.Email, user.Role, s.cfg.JWT.Secret, s.cfg.JWT.ExpireMinutes)
	if err != nil {
		return nil, appErrors.Internal("failed to generate token")
	}

	return &AuthResponse{
		Token:  token,
		APIKey: apiKey,
		User:   *user,
	}, nil
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, appErrors.Internal("failed to query user")
	}
	if user == nil {
		return nil, appErrors.Unauthorized("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, appErrors.Unauthorized("invalid email or password")
	}

	// Auto-assign APIKey if not yet generated
	if user.APIKey == "" {
		user.APIKey = generateAPIKey()
		_ = s.repo.UpdateAPIKey(ctx, user.ID, user.APIKey)
	}

	token, err := middleware.GenerateToken(user.ID, user.Email, user.Role, s.cfg.JWT.Secret, s.cfg.JWT.ExpireMinutes)
	if err != nil {
		return nil, appErrors.Internal("failed to generate token")
	}

	return &AuthResponse{
		Token:  token,
		APIKey: user.APIKey,
		User:   *user,
	}, nil
}

func (s *service) GetProfile(ctx context.Context, userID uint) (*User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, appErrors.Internal("failed to query user")
	}
	if user == nil {
		return nil, appErrors.NotFound("user not found")
	}
	return user, nil
}

func (s *service) RegenerateAPIKey(ctx context.Context, userID uint) (string, error) {
	newKey := generateAPIKey()
	if err := s.repo.UpdateAPIKey(ctx, userID, newKey); err != nil {
		return "", appErrors.Internal("failed to rotate API key")
	}
	return newKey, nil
}

func (s *service) ValidateAPIKey(ctx context.Context, apiKey string) (*User, error) {
	if apiKey == "" {
		return nil, appErrors.Unauthorized("missing API key")
	}

	user, err := s.repo.FindByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, appErrors.Internal("failed to query API key")
	}
	if user == nil {
		return nil, appErrors.Unauthorized("invalid API key")
	}

	return user, nil
}

func (s *service) ValidateAPIKeyForAuth(ctx context.Context, apiKey string) (uint, string, string, error) {
	user, err := s.ValidateAPIKey(ctx, apiKey)
	if err != nil {
		return 0, "", "", err
	}
	return user.ID, user.Email, user.Role, nil
}
