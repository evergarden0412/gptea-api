package main

import (
	"github.com/evergarden0412/gptea-api/internal/server"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	s := server.New()
	r := gin.Default()
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{"https://gptea.keenranger.dev", "https://gptea-test.keenranger.dev"}
	r.Use(cors.New(corsCfg))
	s.Install(r.Handle)
	r.Run(":8080")
}
