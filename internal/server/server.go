package server

import (
	"net/http"
	"os"
	"sort"
	"time"

	_ "github.com/evergarden0412/gptea-api/docs"
	"github.com/evergarden0412/gptea-api/internal"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Server struct {
}

func New() *Server {
	return &Server{}
}

type messageResponse struct {
	Message string `json:"message" example:"Hello, World!"`
}

type errorResponse struct {
	Error string `json:"error" example:"error message"`
}

// @title GPTea API
// @version 0.1.0
// @description This is a sample server for GPTea API.
// @host api.gptea-test.keenranger.dev
// @securityDefinitions.apikey AccessTokenAuth
// @in header
// @name Authorization
// @description type `Bearer {access_token}`
// @securityDefinitions.apikey RefreshTokenAuth
// @in header
// @name X-Refresh-Token
// @description type `Bearer {refresh_token}`
func (s *Server) Install(handle func(string, string, ...gin.HandlerFunc) gin.IRoutes) {
	handle("GET", "/ping2", s.handlePing)
	handle("GET", "/me/chats", s.handleGetMyChats)
	handle("GET", "/me/chats/:chatID/messages", s.handleGetMyMessages)
	handle("GET", "/me/scrapbooks", s.handleGetMyScrapbooks)
	handle("GET", "/me/scrapbooks/:scrapbookID/scraps", s.handleGetMyScraps)
	if os.Getenv("ENV") != "prod" {
		handle("GET", "/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

func (s *Server) handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, messageResponse{Message: "pong"})
}

type chatsResponse struct {
	Chats []internal.Chat `json:"chats"`
}

// handleGetMyChats godoc
// @Summary Get my chats
// @Description Get my chats in descending order of created_at
// @Security AccessTokenAuth
// @Success 200 {object} chatsResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /me/chats [get]
// @tags chats
func (s *Server) handleGetMyChats(c *gin.Context) {
	sampleTime0 := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	sampleTime1 := time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)
	sampleTime2 := time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)
	someSampleChats := []internal.Chat{
		{
			ID:        "1",
			Name:      "chat1",
			CreatedAt: &sampleTime0,
		},
		{
			ID:        "2",
			Name:      "chat2",
			CreatedAt: &sampleTime1,
		},
		{
			ID:        "3",
			Name:      "chat3",
			CreatedAt: &sampleTime2,
		},
	}
	c.JSON(http.StatusOK, chatsResponse{Chats: someSampleChats})
}

// handlePostMyChat godoc
// @Summary Post my chat
// @Description Post my chat
// @Security AccessTokenAuth
// @Success 201 {object} messageResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /me/chats [post]
// @tags chats
func (s *Server) handlePostMyChat(c *gin.Context) {
	c.JSON(http.StatusCreated, messageResponse{Message: "pong"})
}

type messagesResponse struct {
	Messages []internal.Message `json:"messages"`
}

// handleGetMyMessages godoc
// @Summary Get my messages
// @Description Get my messages in descending order of created_at
// @Security AccessTokenAuth
// @Success 200 {object} messagesResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /me/chats/:chatID/messages [get]
// @tags messages
func (s *Server) handleGetMyMessages(c *gin.Context) {
	sampleTime0 := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	sampleTime1 := time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)
	sampleTime2 := time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)
	someSampleMessages := []internal.Message{
		{
			ChatID:    "1",
			Seq:       1,
			Content:   "message1",
			CreatedAt: &sampleTime0,
		},
		{
			ChatID:    "1",
			Seq:       2,
			Content:   "message2",
			CreatedAt: &sampleTime1,
		},
		{
			ChatID:    "1",
			Seq:       3,
			Content:   "message3",
			CreatedAt: &sampleTime2,
		},
	}
	c.JSON(http.StatusOK, messagesResponse{Messages: someSampleMessages})
}

type scrapbooksResponse struct {
	Scrapbooks []internal.Scrapbook `json:"scrapbooks"`
}

// handleGetMyScrapbooks godoc
// @Summary Get my scrapbooks
// @Description Get my scrapbooks
// @Security AccessTokenAuth
// @Success 200 {object} scrapbooksResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /me/scrapbooks [get]
// @tags scraps
func (s *Server) handleGetMyScrapbooks(ctx *gin.Context) {
	sampleTime0 := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	basicScrapbook := internal.Scrapbook{
		ID:        "1",
		Name:      "basic",
		CreatedAt: &sampleTime0,
	}
	ctx.JSON(http.StatusOK, scrapbooksResponse{Scrapbooks: []internal.Scrapbook{basicScrapbook}})
}

type scrapsResponse struct {
	Scraps []internal.Message `json:"scraps"`
}

// handleGetMyScraps godoc
// @Summary Get my scraps
// @Description Get my scraps in descending order of created_at
// @Security AccessTokenAuth
// @Success 200 {object} scrapsResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /me/scrapbooks/:scrapbookID/scraps [get]
// @tags scraps
func (s *Server) handleGetMyScraps(c *gin.Context) {
	sampleTime0 := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	sampleTime1 := time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)
	sampleTime2 := time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)
	someSampleScraps := []internal.Message{
		{
			ChatID:    "1",
			Seq:       1,
			Content:   "message1",
			CreatedAt: &sampleTime0,
		},
		{
			ChatID:    "2",
			Seq:       2,
			Content:   "message2",
			CreatedAt: &sampleTime1,
		},
		{
			ChatID:    "3",
			Seq:       3,
			Content:   "message3",
			CreatedAt: &sampleTime2,
		},
	}
	// order by created_at desc
	sort.Slice(someSampleScraps, func(i, j int) bool {
		return someSampleScraps[i].CreatedAt.After(*someSampleScraps[j].CreatedAt)
	})
	c.JSON(http.StatusOK, scrapsResponse{Scraps: someSampleScraps})
}
