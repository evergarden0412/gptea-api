package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/evergarden0412/gptea-api/internal/auth"
	"github.com/evergarden0412/gptea-api/internal/config"
	"github.com/evergarden0412/gptea-api/internal/postgres"
	"github.com/evergarden0412/gptea-api/internal/server"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kataras/golog"
	_ "github.com/lib/pq"
)

var ginLambda *ginadapter.GinLambda

func main() {
	ctx := context.Background()
	cfg, err := config.Init(ctx)
	accessTokenTTl, err := time.ParseDuration(cfg.AccessTokenTTL)
	if err != nil {
		golog.Fatal(err)
	}
	refreshTokenTTl, err := time.ParseDuration(cfg.RefreshTokenTTL)
	if err != nil {
		golog.Fatal(err)
	}
	a := auth.New(auth.AuthenticatorConfig{
		AccessTokenTTL:  accessTokenTTl,
		RefreshTokenTTL: refreshTokenTTl,
		AccessTokenKey:  []byte(cfg.AccessTokenKey),
		RefreshTokenKey: []byte(cfg.RefreshTokenKey),
	})
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, "gptea"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	postgresDB := postgres.New(db)
	s := server.New(a, postgresDB)
	r := gin.Default()
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{"https://gptea.keenranger.dev", "https://gptea-test.keenranger.dev"}
	corsCfg.AllowHeaders = []string{"origin", "content-length", "content-type", "authorization", "x-refresh-token"}
	r.Use(cors.New(corsCfg))
	s.Install(r.Handle)
	if os.Getenv("LOCAL") == "true" {
		r.Run(":8080")
		return
	}
	ginLambda = ginadapter.New(r)
	lambda.Start(Handler)
}
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}
