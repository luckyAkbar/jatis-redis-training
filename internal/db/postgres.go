package db

import (
	"os"

	"github.com/luckyAkbar/jatis-redis-training/internal/config"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	// PostgresDB postgres database connection
	PostgresDB *gorm.DB
)

// InitializePostgresConn initialize postgres connection
// call os.Exit() if error accoured
func InitializePostgresConn() {
	conn, err := initPostgresConn(config.PostgresDSN())
	if err != nil {
		logrus.Error("failed to connect to postgres database. reason: ", err.Error())
		os.Exit(1)
	}

	PostgresDB = conn

	switch config.LogLevel() {
	case "error":
		PostgresDB.Logger = PostgresDB.Logger.LogMode(gormLogger.Error)
	case "warn":
		PostgresDB.Logger = PostgresDB.Logger.LogMode(gormLogger.Warn)
	default:
		PostgresDB.Logger = PostgresDB.Logger.LogMode(gormLogger.Info)
	}

	logrus.Info("Connected to Postgres Database")
}

func initPostgresConn(dsn string) (*gorm.DB, error) {
	conn, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}

	return conn, err
}
