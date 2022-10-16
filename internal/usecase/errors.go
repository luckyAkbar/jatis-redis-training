package usecase

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrBadRequest    = errors.New("bad request")
	ErrInternal      = errors.New("internal error")
	ErrNotFound      = errors.New("not found")
	ErrUnauthorized  = errors.New("unauthorized request")
	ErrAlreadyExists = errors.New("already exists")

	HTTPErrUnauthorized = echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
)
