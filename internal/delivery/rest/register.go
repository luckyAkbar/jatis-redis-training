package rest

import (
	"net/http"

	"github.com/kumparan/go-utils"
	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/jatis-redis-training/internal/model"
	"github.com/luckyAkbar/jatis-redis-training/internal/usecase"
	"github.com/sirupsen/logrus"
)

func (s *RESTService) handleRegister() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := &model.RegisterInput{}
		if err := ctx.Bind(input); err != nil {
			return ErrBadRequest
		}

		user, err := s.registerUsecase.Register(ctx.Request().Context(), input)
		switch err {
		default:
			logrus.WithFields(logrus.Fields{
				"ctx":   utils.DumpIncomingContext(ctx.Request().Context()),
				"input": utils.Dump(input),
			}).Error(err)

			return ErrInternal
		case usecase.ErrBadRequest:
			return ErrBadRequest
		case usecase.ErrAlreadyExists:
			return ErrAlreadyExists
		case nil:
			return ctx.JSON(http.StatusCreated, user)
		}
	}
}
