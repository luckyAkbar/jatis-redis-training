package rest

import (
	"net/http"

	"github.com/kumparan/go-utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *RESTService) handleGetMenu() echo.HandlerFunc {
	return func(c echo.Context) error {
		menu, err := s.contentUsecase.GetAllMenu(c.Request().Context())
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"ctx": utils.DumpIncomingContext(c.Request().Context()),
			}).Error(err)
		}

		return c.JSON(http.StatusOK, menu)
	}
}

func (s *RESTService) handleGetData() echo.HandlerFunc {
	return func(c echo.Context) error {
		data, err := s.contentUsecase.GetAllData(c.Request().Context())
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"ctx": utils.DumpIncomingContext(c.Request().Context()),
			}).Error(err)
		}

		return c.JSON(http.StatusOK, data)
	}
}
