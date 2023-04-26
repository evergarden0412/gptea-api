package server

import (
	"net/http"
	"os"

	_ "github.com/evergarden0412/gptea-api/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Server struct {
}

func New() *Server {
	return &Server{}
}

type MessageResponse struct {
	Message string `json:"message" example:"Hello, World!"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

// @title GPTea API
// @version 0.1.0
// @description This is a sample server for GPTea API.
// @host api.gptea-test.keenranger.dev
func (s *Server) Install(handle func(string, string, ...gin.HandlerFunc) gin.IRoutes) {
	handle("GET", "/ping2", s.handlePing)
	handle("GET", "/me/chats", s.handleGetMyChats)
	handle("GET", "/me/scraps", s.handleGetMyScraps)
	handle("GET", "/me/highlights", s.handleGetMyHighlights)
	handle("GET", "/chats/:chatID/messages", s.handleGetMessages)
	if os.Getenv("ENV") != "prod" {
		handle("GET", "/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

func (s *Server) handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, MessageResponse{Message: "pong"})
}

func (s *Server) handleGetMyChats(c *gin.Context) {
	c.JSON(http.StatusOK, MessageResponse{Message: "pong"})
}

func (s *Server) handleGetMyScraps(c *gin.Context) {
	c.JSON(http.StatusOK, MessageResponse{Message: "pong"})
}

func (s *Server) handleGetMyHighlights(c *gin.Context) {
	c.JSON(http.StatusOK, MessageResponse{Message: "pong"})
}

func (s *Server) handleGetMessages(c *gin.Context) {
	c.JSON(http.StatusOK, MessageResponse{Message: "pong"})
}
