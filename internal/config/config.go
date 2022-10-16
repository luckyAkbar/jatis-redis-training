package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var (
	DefaultTokenLength     = 25
	DefaultSessionExpiry   = 24 * time.Hour * 30
	DefaultDataCacheExpiry = 24 * time.Hour
	DefaultMenuCacheExpiry = 24 * time.Hour
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")

	err := viper.ReadInConfig()
	if err != nil {
		panic("failed to read config file")
	}
}

func Env() string {
	env := viper.GetString("env")
	if env == "" {
		return "development"
	}

	return env
}

func LogLevel() string {
	level := viper.GetString("server.log.level")
	if level == "" {
		return "DEBUG"
	}

	return level
}

func ServerPort() string {
	port := viper.GetString("server.port")
	if port == "" {
		return ":8080"
	}

	return fmt.Sprintf(":%s", port)
}

func PostgresDSN() string {
	host := viper.GetString("postgres.host")
	db := viper.GetString("postgres.db")
	user := viper.GetString("postgres.user")
	pw := viper.GetString("postgres.password")
	port := viper.GetString("postgres.port")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, pw, db, port)
}

func RedisAddr() string {
	return viper.GetString("redis.addr")
}

func RedisPassword() string {
	return viper.GetString("redis.password")
}

func RedisCacheDB() int {
	return viper.GetInt("redis.db")
}

func RedisMinIdleConn() int {
	return viper.GetInt("redis.min")
}

func RedisMaxIdleConn() int {
	return viper.GetInt("redis.max")
}
