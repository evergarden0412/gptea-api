package config

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kataras/golog"
)

type Config struct {
	Env             string
	Region          string
	AccessTokenTTL  string
	RefreshTokenTTL string
	AccessTokenKey  string
	RefreshTokenKey string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
}

type PostgresSecret struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"username"`
	Password string `json:"password"`
}

func Init(ctx context.Context) (*Config, error) {
	golog.SetLevel("debug")
	if os.Getenv("ENV") == "prod" {
		golog.SetLevel("error")
		gin.SetMode(gin.ReleaseMode)
	}
	cfg := &Config{}
	cfg.Env = os.Getenv("ENV")
	cfg.Region = os.Getenv("REGION")
	if cfg.Region == "" {
		cfg.Region = "ap-northeast-2"
	}
	if os.Getenv("LOCAL") == "true" {
		cfg.AccessTokenTTL = "5m"
		cfg.RefreshTokenTTL = "10m"
		cfg.AccessTokenKey = "access_token_secret"
		cfg.RefreshTokenKey = "refresh_token_secret"
		cfg.DBHost = "localhost"
		cfg.DBPort = "5432"
		cfg.DBUser = "postgres"
		cfg.DBPassword = "password"
		return cfg, nil
	}
	return cfg, nil
}
