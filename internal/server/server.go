package server

import (
	"net/http"
	"os"
	"sort"
	"time"

	_ "github.com/evergarden0412/gptea-api/docs"
	"github.com/evergarden0412/gptea-api/internal"
	"github.com/evergarden0412/gptea-api/internal/auth"
	"github.com/evergarden0412/gptea-api/internal/chatbot"
	"github.com/evergarden0412/gptea-api/internal/postgres"
	"github.com/gin-gonic/gin"
	"github.com/kataras/golog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Server struct {
	c  *chatbot.Chatbot
	a  *auth.Authenticator
	db *postgres.DB
}

func New(a *auth.Authenticator, chatbot *chatbot.Chatbot, db *postgres.DB) *Server {
	return &Server{
		a:  a,
		c:  chatbot,
		db: db,
	}
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
// @name authorization
// @description type `Bearer {access_token}`
// @securityDefinitions.apikey RefreshTokenAuth
// @in header
// @name x-refresh-token
// @description type `{refresh_token}`
func (s *Server) Install(handle func(string, string, ...gin.HandlerFunc) gin.IRoutes) {
	handle("GET", "/ping2", s.handlePing)
	handle("POST", "/auth/cred/oauth/token")
	handle("POST", "/auth/cred/register", s.handleRegister)
	handle("POST", "/auth/cred/sign-in", s.handleSignIn)
	handle("GET", "/auth/token/verify", s.handleVerifyToken)
	handle("POST", "/auth/token/refresh", s.handleRefreshToken)
	handle("DELETE", "/me", s.ensureUser, s.handleDeleteMe)
	// chat
	handle("GET", "/me/chats", s.ensureUser, s.handleGetMyChats)
	handle("POST", "/me/chats", s.ensureUser, s.handlePostMyChat)
	handle("PATCH", "/me/chats/:chatID", s.ensureUser, s.handlePatchMyChat)
	handle("DELETE", "/me/chats/:chatID", s.ensureUser, s.handleDeleteMyChat)
	// message
	handle("GET", "/me/chats/:chatID/messages", s.ensureUser, s.handleGetMyMessages)
	handle("POST", "/me/chats/:chatID/messages", s.ensureUser, s.handlePostMyMessage)
	handle("GET", "/me/scrapbooks", s.ensureUser, s.handleGetMyScrapbooks)
	handle("GET", "/me/scrapbooks/:scrapbookID/scraps", s.ensureUser, s.handleGetMyScraps)
	handle("DELETE", "/me/scrapbooks/", s.ensureUser)
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

// handleDeleteMe godoc
// @Summary Delete my account
// @Description Delete my account
// @Security AccessTokenAuth
// @Success 204
// @Failure 500 {object} errorResponse
// @Router /me [delete]
// @tags users
func (s *Server) handleDeleteMe(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	if err := s.db.Resign(ctx, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		golog.Error("handleDeleteMe: delete user: ", err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

type messagesResponse struct {
	Messages []internal.Message `json:"messages"`
}

// handleGetMyMessages godoc
// @summary Get my messages
// @description Get my messages in descending order of created_at
// @tags messages
// @security AccessTokenAuth
// @param chatID path string true "chatID"
// @success 200 {object} messagesResponse
// @failure 400 {object} errorResponse
// @failure 500 {object} errorResponse
// @router /me/chats/{chatID}/messages [get]
func (s *Server) handleGetMyMessages(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	chatID := ctx.Param("chatID")

	messages, err := s.db.GetMyMessages(ctx, userID, chatID)
	if err != nil {
		golog.Error("handleGetMyMessages: get messages: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	messagesResp := make([]internal.Message, len(messages))
	for i, message := range messages {
		messagesResp[i] = *message
	}

	ctx.JSON(http.StatusOK, messagesResponse{Messages: messagesResp})
}

type messageBody struct {
	Content string `json:"content"`
}

// handlePostMyMessage godoc
// @summary Post my message
// @description Post my message and get response when chatbot finishes processing
// @tags messages
// @security AccessTokenAuth
// @param chatID path string true "chatID"
// @param body body messageBody true "body"
// @success 201 {object} messageResponse
// @failure 400 {object} errorResponse
// @failure 500 {object} errorResponse
// @router /me/chats/{chatID}/messages [post]
func (s *Server) handlePostMyMessage(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	chatID := ctx.Param("chatID")
	var body messageBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		golog.Error("handlePostMyMessage: bind json: ", err)
		return
	}

	history, err := s.db.GetMyMessages(ctx, userID, chatID)
	if err != nil {
		golog.Error("handlePostMyMessage: get history: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	inMsg, outMsg, err := s.c.SendChat(ctx, chatID, history, body.Content)
	if err != nil {
		golog.Error("handlePostMyMessage: send chat: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	if err := s.db.InsertMessage(ctx, userID, *inMsg); err != nil {
		golog.Error("handlePostMyMessage: insert message: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	if err := s.db.InsertMessage(ctx, userID, *outMsg); err != nil {
		golog.Error("handlePostMyMessage: insert message: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, messageResponse{Message: outMsg.Content})
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
// @Param scrapbookID path string true "scrapbookID"
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
