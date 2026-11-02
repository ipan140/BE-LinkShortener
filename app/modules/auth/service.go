package auth

import (
	"context"
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

	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, appErrors.Internal("failed to create user")
	}

	token, err := middleware.GenerateToken(user.ID, user.Email, user.Role, s.cfg.JWT.Secret, s.cfg.JWT.ExpireMinutes)
	if err != nil {
		return nil, appErrors.Internal("failed to generate token")
	}

	return &AuthResponse{
		Token: token,
		User:  *user,
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

	token, err := middleware.GenerateToken(user.ID, user.Email, user.Role, s.cfg.JWT.Secret, s.cfg.JWT.ExpireMinutes)
	if err != nil {
		return nil, appErrors.Internal("failed to generate token")
	}

	return &AuthResponse{
		Token: token,
		User:  *user,
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
