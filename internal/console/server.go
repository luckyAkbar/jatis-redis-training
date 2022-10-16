package console

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/luckyAkbar/jatis-redis-training/internal/config"
	"github.com/luckyAkbar/jatis-redis-training/internal/db"
	"github.com/luckyAkbar/jatis-redis-training/internal/delivery/rest"
	"github.com/luckyAkbar/jatis-redis-training/internal/repository"
	"github.com/luckyAkbar/jatis-redis-training/internal/usecase"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start the server",
	Run:   server,
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

func server(c *cobra.Command, args []string) {
	db.InitializePostgresConn()

	sqlDB, err := db.PostgresDB.DB()
	if err != nil {
		logrus.Fatal("unable to start server. reason: ", err.Error())
	}

	defer sqlDB.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:         config.RedisAddr(),
		Password:     config.RedisPassword(),
		DB:           config.RedisCacheDB(),
		MinIdleConns: config.RedisMinIdleConn(),
		MaxIdleConns: config.RedisMaxIdleConn(),
	})

	_, err = redisClient.Ping(context.TODO()).Result()
	if err != nil {
		logrus.Error(fmt.Sprintf("failed to connect to redis: %v", err))
	}

	logrus.Info("Connected to redis server")

	cacher := db.NewCacher(redisClient)

	userRepo := repository.NewUserRepo(db.PostgresDB)
	sessionRepo := repository.NewSessionRepo(db.PostgresDB, cacher)
	contentRepo := repository.NewContentRepository(db.PostgresDB, cacher)

	registerUsecase := usecase.NewRegisterUsecase(userRepo)
	authUsecase := usecase.NewAuthUsecase(sessionRepo, userRepo)
	contentUsecase := usecase.NewContentUsecase(contentRepo)

	HTTPServer := echo.New()

	HTTPServer.Pre(middleware.AddTrailingSlash())
	HTTPServer.Use(middleware.Logger())
	HTTPServer.Use(middleware.Recover())
	HTTPServer.Use(middleware.CORS())
	HTTPServer.Use(authUsecase.CreateAuthMiddleware())
	HTTPServer.Use(authUsecase.CreateRejectUnauthorizedRequestMiddleware([]string{
		"/soal_redis/login/",
		"/soal_redis/register/",
	}))

	apiGroup := HTTPServer.Group("soal_redis")
	rest.InitRESTService(apiGroup, registerUsecase, authUsecase, contentUsecase)

	if err := HTTPServer.Start(config.ServerPort()); err != nil {
		log.Fatal("unable to start server: ", err)
	}
}
