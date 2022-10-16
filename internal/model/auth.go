package model

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

type AuthUser struct {
	Username string `json:"username"`
}

type Session struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	AccessToken string    `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiredAt   time.Time `json:"expired_at"`
}

// IsAccessTokenExpired compare session access token against time.Now()
// return true if the token is expired by now
func (s *Session) IsAccessTokenExpired() bool {
	if s == nil {
		return true
	}

	now := time.Now()
	return now.After(s.ExpiredAt)
}

type AuthUsecase interface {
	Login(ctx context.Context, input *LoginInput) (*Session, error)
	CreateAuthMiddleware() echo.MiddlewareFunc
	CreateRejectUnauthorizedRequestMiddleware(skipURL []string) echo.MiddlewareFunc
}

type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	FindByAccessToken(ctx context.Context, accessToken string) (*Session, error)
}
