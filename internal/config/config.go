package config

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
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
	OpenAIAPIKey    string
	OpenAIAPIOrgID  string
}

type PostgresSecret struct {
	User     string `json:"username"`
	Password string `json:"password"`
}

type HMACSecret struct {
	AccessTokenKey  string `json:"accessTokenKey"`
	RefreshTokenKey string `json:"refreshTokenKey"`
}

type OpenAIAPISecret struct {
	Key            string `json:"key"`
	OrganizationID string `json:"organizationID"`
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
	cfg.AccessTokenTTL = os.Getenv("ACCESS_TOKEN_TTL")
	cfg.RefreshTokenTTL = os.Getenv("REFRESH_TOKEN_TTL")
	cfg.DBHost = os.Getenv("DB_HOST")
	cfg.DBPort = os.Getenv("DB_PORT")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region),
	})
	secretsManager := secretsmanager.New(sess)

	postgresInput := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(os.Getenv("DB_SECRET_ARN")),
	}
	postgresOutput, err := secretsManager.GetSecretValue(postgresInput)
	if err != nil {
		return nil, err
	}
	var postgresString string
	if postgresOutput.SecretString != nil {
		postgresString = *postgresOutput.SecretString
	} else {
		postgresString = string(postgresOutput.SecretBinary)
	}
	var postgresSecret PostgresSecret
	if err := json.Unmarshal([]byte(postgresString), &postgresSecret); err != nil {
		return nil, err
	}

	hmacInput := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(os.Getenv("HMAC_SECRET_ARN")),
	}
	hmacOutput, err := secretsManager.GetSecretValue(hmacInput)
	if err != nil {
		return nil, err
	}
	var hmacString string
	if hmacOutput.SecretString != nil {
		hmacString = *hmacOutput.SecretString
	} else {
		hmacString = string(hmacOutput.SecretBinary)
	}
	var hmacSecret HMACSecret
	if err := json.Unmarshal([]byte(hmacString), &hmacSecret); err != nil {
		return nil, err
	}

	openaiAPIKeyInput := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(os.Getenv("OPENAI_API_SECRET_ARN")),
	}
	openaiAPIKeyOutput, err := secretsManager.GetSecretValue(openaiAPIKeyInput)
	if err != nil {
		return nil, err
	}
	var openAIAPIString string
	if openaiAPIKeyOutput.SecretString != nil {
		openAIAPIString = *openaiAPIKeyOutput.SecretString
	} else {
		openAIAPIString = string(openaiAPIKeyOutput.SecretBinary)
	}
	var openAIAPISecret OpenAIAPISecret
	if err := json.Unmarshal([]byte(openAIAPIString), &openAIAPISecret); err != nil {
		return nil, err
	}

	cfg.DBUser = postgresSecret.User
	cfg.DBPassword = postgresSecret.Password
	cfg.AccessTokenKey = hmacSecret.AccessTokenKey
	cfg.RefreshTokenKey = hmacSecret.RefreshTokenKey
	cfg.OpenAIAPIKey = openAIAPISecret.Key
	cfg.OpenAIAPIOrgID = openAIAPISecret.OrganizationID
	return cfg, nil
}
