package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrBadRequest    = echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	ErrInternal      = echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	ErrUnauthorized  = echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	ErrNotFound      = echo.NewHTTPError(http.StatusNotFound, "Not Found")
	ErrAlreadyExists = echo.NewHTTPError(http.StatusForbidden, "Already Exists")
)
