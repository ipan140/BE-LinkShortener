package middleware

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	appErrors "link-shortener/common/errors"
	"link-shortener/common/utils"
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type APIKeyValidator interface {
	ValidateAPIKey(ctx context.Context, apiKey string) (any, error)
}

// GenerateToken creates a signed JWT token
func GenerateToken(userID uint, email, role, secret string, expireMinutes int) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

type UserWithDetails interface {
	GetID() uint
	GetEmail() string
	GetRole() string
}

type KeyFinder interface {
	ValidateAPIKey(ctx context.Context, key string) (uint, string, string, error)
}

// CombinedAuthMiddleware allows authentication using either 'X-API-Key' or 'Authorization: Bearer <token>'
func CombinedAuthMiddleware(secret string, finder any) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 1. Check for X-API-Key header first
			apiKey := strings.TrimSpace(c.Request().Header.Get("X-API-Key"))
			if apiKey != "" {
				if kf, ok := finder.(interface {
					ValidateAPIKey(ctx context.Context, apiKey string) (any, error)
				}); ok {
					user, err := kf.ValidateAPIKey(c.Request().Context(), apiKey)
					if err == nil && user != nil {
						if u, ok := user.(interface {
							GetAuthData() (uint, string, string)
						}); ok {
							uid, email, role := u.GetAuthData()
							c.Set("user_id", uid)
							c.Set("email", email)
							c.Set("role", role)
							return next(c)
						}
					}
				}
				// Also handle direct concrete struct via reflection or interface
				if kf, ok := finder.(interface {
					ValidateAPIKeyForAuth(ctx context.Context, apiKey string) (uint, string, string, error)
				}); ok {
					uid, email, role, err := kf.ValidateAPIKeyForAuth(c.Request().Context(), apiKey)
					if err == nil {
						c.Set("user_id", uid)
						c.Set("email", email)
						c.Set("role", role)
						return next(c)
					}
				}
			}

			// 2. Fallback to JWT Bearer Token
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return utils.JSONError(c, appErrors.Unauthorized("missing Authorization header or X-API-Key"))
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return utils.JSONError(c, appErrors.Unauthorized("invalid Authorization header format"))
			}

			tokenString := parts[1]
			token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				return utils.JSONError(c, appErrors.Unauthorized("invalid or expired token"))
			}

			claims, ok := token.Claims.(*JWTClaims)
			if !ok {
				return utils.JSONError(c, appErrors.Unauthorized("failed to parse token claims"))
			}

			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}

// JWTMiddleware validates the bearer token and injects claims into echo context
func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return CombinedAuthMiddleware(secret, nil)
}
