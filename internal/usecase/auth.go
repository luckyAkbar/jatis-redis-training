package usecase

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/kumparan/go-utils"
	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/jatis-redis-training/internal/config"
	"github.com/luckyAkbar/jatis-redis-training/internal/helper"
	"github.com/luckyAkbar/jatis-redis-training/internal/model"
	"github.com/luckyAkbar/jatis-redis-training/internal/repository"
	"github.com/sirupsen/logrus"
)

type contextKey string

const (
	userCtxKey           contextKey = "github.com/luckyAkbar/jatis-redis-training/internal/auth.User"
	_headerAuthorization string     = "Authorization"
	_authScheme          string     = "Bearer"
)

type authUsecase struct {
	sessionRepo model.SessionRepository
	userRepo    model.UserRepository
}

func NewAuthUsecase(sessionRepo model.SessionRepository, userRepo model.UserRepository) model.AuthUsecase {
	return &authUsecase{
		sessionRepo,
		userRepo,
	}
}

func (u *authUsecase) Login(ctx context.Context, input *model.LoginInput) (*model.Session, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"input": utils.Dump(input),
	})

	if err := input.Validate(); err != nil {
		return nil, ErrBadRequest
	}

	user, err := u.userRepo.GetByUsername(ctx, input.Username)
	switch err {
	default:
		logger.Error(err)
		return nil, ErrInternal
	case repository.ErrNotFound:
		return nil, ErrNotFound
	case nil:
		break
	}

	if user.Password != input.Password {
		return nil, ErrUnauthorized
	}

	session := &model.Session{
		ID:          utils.GenerateID(),
		Username:    user.Username,
		AccessToken: helper.GenerateToken(utils.GenerateID()),
		CreatedAt:   time.Now(),
		ExpiredAt:   time.Now().Add(config.DefaultSessionExpiry),
	}

	if err := u.sessionRepo.Create(ctx, session); err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	return session, nil
}

func (u *authUsecase) CreateAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := u.getTokenFromHeader(c.Request())
			if token == "" {
				// pass the auth process to next handler
				return next(c)
			}

			ctx := c.Request().Context()

			session, err := u.getSessionFromAccessToken(ctx, token)
			if err != nil {
				return next(c)
			}

			// just pass if expired. The next middleware should block the request
			// if needed
			if session.IsAccessTokenExpired() {
				return next(c)
			}

			user, err := u.userRepo.GetByUsername(ctx, session.Username)
			if err != nil {
				return next(c)
			}

			ctx = SetUserToCtx(ctx, model.AuthUser{
				Username: user.Username,
			})

			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func (u *authUsecase) CreateRejectUnauthorizedRequestMiddleware(skipURL []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			url := c.Request().URL
			for _, skip := range skipURL {
				if url.Path == skip {
					return next(c)
				}
			}

			user := GetUserFromCtx(c.Request().Context())
			if user == nil {
				return HTTPErrUnauthorized
			}

			return next(c)
		}
	}
}

func (u *authUsecase) getTokenFromHeader(req *http.Request) string {
	authHeader := strings.Split(req.Header.Get(_headerAuthorization), " ")

	if len(authHeader) != 2 || authHeader[0] != _authScheme {
		return ""
	}

	return strings.TrimSpace(authHeader[1])
}

func (u *authUsecase) getSessionFromAccessToken(ctx context.Context, token string) (*model.Session, error) {
	logger := logrus.WithField("token", token)

	session, err := u.sessionRepo.FindByAccessToken(ctx, token)
	if err != nil {
		logger.Warn(err)
		return nil, err
	}

	return session, nil
}

// SetUserToCtx self explained
func SetUserToCtx(ctx context.Context, user model.AuthUser) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

// GetUserFromCtx self explained
func GetUserFromCtx(ctx context.Context) *model.AuthUser {
	user, ok := ctx.Value(userCtxKey).(model.AuthUser)
	if !ok {
		return nil
	}

	return &user
}
