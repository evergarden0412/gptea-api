package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/evergarden0412/gptea-api/internal/server"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func main() {
	s := server.New()
	r := gin.Default()
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{"https://gptea.keenranger.dev", "https://gptea-test.keenranger.dev"}
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
