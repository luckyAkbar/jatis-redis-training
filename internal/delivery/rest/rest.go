package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/jatis-redis-training/internal/model"
)

type RESTService struct {
	group           *echo.Group
	registerUsecase model.RegisterUsecase
	authUsecase     model.AuthUsecase
	contentUsecase  model.ContentUsecase
}

func InitRESTService(group *echo.Group, registerUsecase model.RegisterUsecase, authUsecase model.AuthUsecase, contentUsecase model.ContentUsecase) {
	service := &RESTService{
		group,
		registerUsecase,
		authUsecase,
		contentUsecase,
	}

	service.initRoutes()
}

func (s *RESTService) initRoutes() {
	s.group.POST("/register/", s.handleRegister())
	s.group.POST("/login/", s.handleLogin())
	s.group.GET("/data/", s.handleGetData())
	s.group.GET("/menu/", s.handleGetMenu())
}
