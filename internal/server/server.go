package server

import (
	"net/http"
	"os"

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
	handle("POST", "/auth/cred/register", s.handleRegister)
	handle("POST", "/auth/cred/sign-in", s.handleSignIn)
	handle("POST", "/auth/cred/logout", s.ensureUser, s.handleLogout)
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
	// scrapbook
	handle("GET", "/me/scrapbooks", s.ensureUser, s.handleGetMyScrapbooks)
	handle("POST", "/me/scrapbooks", s.ensureUser, s.handlePostMyScrapbook)
	handle("DELETE", "/me/scrapbooks/:scrapbookID", s.ensureUser, s.handleDeleteMyScrapbook)
	handle("PATCH", "/me/scrapbooks/:scrapbookID", s.ensureUser, s.handlePatchMyScrapbook)
	// scrap
	handle("GET", "/me/scrapbooks/:scrapbookID/scraps", s.ensureUser, s.handleGetScrapsOnScrapbook)
	handle("GET", "/me/scraps", s.ensureUser, s.handleGetMyScraps)
	handle("POST", "/me/scraps", s.ensureUser, s.handlePostMyScrap)
	handle("DELETE", "/me/scraps/:scrapID", s.ensureUser, s.handleDeleteMyScrap)
	handle("GET", "/me/scraps/:scrapID/scrapbooks", s.ensureUser, s.handleGetMyScrapbooksOnScrap)
	handle("POST", "/me/scraps/:scrapID/scrapbooks/:scrapbookID", s.ensureUser, s.handlePostScrapOnScrapbook)
	handle("DELETE", "/me/scraps/:scrapID/scrapbooks/:scrapbookID", s.ensureUser, s.handleDeleteScrapOnScrapbook)
	if os.Getenv("ENV") != "prod" {
		handle("GET", "/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

func (s *Server) handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, messageResponse{Message: "pong"})
}

type messagesResponse struct {
	Messages []internal.MessageWithScrap `json:"messages"`
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

	messagesResp := make([]internal.MessageWithScrap, len(messages))
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
//
//	@summary Get my scrapbooks
//	@description Get my scrapbooks
//	@tags scrapbooks
//	@security AccessTokenAuth
//	@success 200 {object} scrapbooksResponse
//	@failure 400 {object} errorResponse
//	@failure 500 {object} errorResponse
//	@router /me/scrapbooks [get]
func (s *Server) handleGetMyScrapbooks(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	scrapbooks, err := s.db.SelectMyScrapbooks(ctx, userID)
	if err != nil {
		golog.Error("handleGetMyScrapbooks: select my scrapbooks: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	scrapbooksResp := make([]internal.Scrapbook, len(scrapbooks))
	for i, scrapbook := range scrapbooks {
		scrapbooksResp[i] = scrapbook
	}
	ctx.JSON(http.StatusOK, scrapbooksResponse{Scrapbooks: scrapbooksResp})
}

type scrapbookBody struct {
	Name string `json:"name"`
}

// handlePostMyScrapbook godoc
//
//	@summary post new scrapbook
//	@description post new scrapbook
//	@tags scrapbooks
//	@security AccessTokenAuth
//	@param body body scrapbookBody true "body"
//	@success 201
//	@failure 400 {object} errorResponse
//	@failure 500 {object} errorResponse
//	@router /me/scrapbooks [post]
func (s *Server) handlePostMyScrapbook(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	var body scrapbookBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		golog.Error("handlePostMyScrapbook: bind json: ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	scrapbook := internal.Scrapbook{
		Name: body.Name,
	}
	if err := scrapbook.Assign(); err != nil {
		golog.Error("handlePostMyScrapbook: assign: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	if err := s.db.InsertScrapbook(ctx, userID, scrapbook); err != nil {
		golog.Error("handlePostMyScrapbook: insert scrapbook: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusCreated)
}

// handleDeleteMyScrapbook godoc
//
//	@summary delete scrapbook
//	@description delete scrapbook
//	@tags scrapbooks
//	@security AccessTokenAuth
//	@param scrapbookID path string true "scrapbookID"
//	@success 204
//	@failure 400 {object} errorResponse
//	@failure 500 {object} errorResponse
//	@router /me/scrapbooks/{scrapbookID} [delete]
func (s *Server) handleDeleteMyScrapbook(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	scrapbookID := ctx.Param("scrapbookID")

	if err := s.db.DeleteScrapbook(ctx, userID, scrapbookID); err != nil {
		golog.Error("handleDeleteMyScrapbook: delete scrapbook: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// handlePatchMyScrapbook godoc
//
//	@summary patch scrapbook
//	@description patch scrapbook
//	@tags scrapbooks
//	@security AccessTokenAuth
//	@param scrapbookID path string true "scrapbookID"
//	@param body body scrapbookBody true "body"
//	@success 204
//	@failure 400 {object} errorResponse
//	@failure 500 {object} errorResponse
//	@router /me/scrapbooks/{scrapbookID} [patch]
func (s *Server) handlePatchMyScrapbook(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	scrapbookID := ctx.Param("scrapbookID")
	var body scrapbookBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		golog.Error("handlePatchMyScrapbook: bind json: ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	if err := s.db.PatchScrapbook(ctx, userID, scrapbookID, body.Name); err != nil {
		golog.Error("handlePatchMyScrapbook: patch scrapbook: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

type scrapsResponse struct {
	Scraps []internal.ScrapWithMessage `json:"scraps"`
}

// handleGetScrapsOnScrapbook godoc
//
//	@summary get scraps on scrapbook
//	@description get scraps on scrapbook
//	@tags scraps
//	@security AccessTokenAuth
//	@param scrapbookID path string true "scrapbookID"
//	@success 200 {object} scrapsResponse
//	@failure 400 {object} errorResponse
//	@failure 500 {object} errorResponse
//	@router /me/scrapbooks/{scrapbookID}/scraps [get]
func (s *Server) handleGetScrapsOnScrapbook(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	scrapbookID := ctx.Param("scrapbookID")

	scraps, err := s.db.SelectScrapsOnScrapbook(ctx, userID, scrapbookID)
	if err != nil {
		golog.Error("handleGetScrapsOnScrapbook: select scraps on scrapbook: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	scrapsResp := make([]internal.ScrapWithMessage, len(scraps))
	for i, scrap := range scraps {
		scrapsResp[i] = scrap
	}
	ctx.JSON(http.StatusOK, scrapsResponse{Scraps: scrapsResp})
}

type postScrapBody struct {
	ChatID       string   `json:"chatID" binding:"required"`
	Seq          int      `json:"seq" binding:"required"`
	Memo         string   `json:"memo"`
	ScrapbookIDs []string `json:"scrapbookIDs" binding:"required,min=1"`
}

// handleGetMyScraps godoc
//
//	@summary get my scraps
//	@description get all my scraps
//	@tags scraps
//	@security AccessTokenAuth
//	@success 200 {object} scrapsResponse
//	@failure 400 {object} errorResponse
//	@failure 500 {object} errorResponse
//	@router /me/scraps [get]
func (s *Server) handleGetMyScraps(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	scraps, err := s.db.SelectMyScraps(ctx, userID)
	if err != nil {
		golog.Error("handleGetMyScraps: select my scraps: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	scrapsResp := make([]internal.ScrapWithMessage, len(scraps))
	for i, scrap := range scraps {
		scrapsResp[i] = scrap
	}
	ctx.JSON(http.StatusOK, scrapsResponse{Scraps: scrapsResp})
}

// handlePostMyScrap godoc
//
//	@summary post new scrap
//	@description post new scrap, store in default scrapbook
//	@tags scraps
//	@security AccessTokenAuth
//	@param body body postScrapBody true "body"
//	@success 201
//	@failure 400 {object} errorResponse
//	@failure 500 {object} errorResponse
//	@router /me/scraps [post]
func (s *Server) handlePostMyScrap(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	var body postScrapBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		golog.Error("handlePostMyScrap: bind json: ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	scrap := internal.Scrap{
		Memo: body.Memo,
	}
	msg := internal.Message{
		ChatID: body.ChatID,
		Seq:    body.Seq,
	}
	if err := scrap.Assign(); err != nil {
		golog.Error("handlePostMyScrap: assign: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	if err := s.db.InsertScrap(ctx, userID, scrap, msg, body.ScrapbookIDs); err != nil {
		golog.Error("handlePostMyScrap: insert scrap: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusCreated)
}

// handleDeleteMyScrap godoc
//
//	@summary delete scrap
//	@description delete scrap
//	@tags scraps
//	@security AccessTokenAuth
//	@param scrapID path string true "scrapID"
//	@success 204
//	@failure 400 {object} errorResponse
//	@failure 500 {object} errorResponse
//	@router /me/scraps/{scrapID} [delete]
func (s *Server) handleDeleteMyScrap(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	scrapID := ctx.Param("scrapID")

	if err := s.db.DeleteScrap(ctx, userID, scrapID); err != nil {
		golog.Error("handleDeleteMyScrap: delete scrap: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// handleGetMyScrapbooksOnScrap godoc
//
//	@summary get my scrapbooks on scrap
//	@description get my scrapbooks on scrap
//	@tags scraps
//	@security AccessTokenAuth
//	@param scrapID path string true "scrapID"
//	@success 200 {object} scrapbooksResponse
//	@failure 400 {object} errorResponse
//	@failure 500 {object} errorResponse
//	@router /me/scraps/{scrapID}/scrapbooks [get]
func (s *Server) handleGetMyScrapbooksOnScrap(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	scrapID := ctx.Param("scrapID")

	scrapbooks, err := s.db.SelectMyScrapbooksOnScrap(ctx, userID, scrapID)
	if err != nil {
		golog.Error("handleGetMyScrapbooksOnScrap: select my scrapbooks on scrap: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	scrapbooksResp := make([]internal.Scrapbook, len(scrapbooks))
	for i, scrapbook := range scrapbooks {
		scrapbooksResp[i] = scrapbook
	}
	ctx.JSON(http.StatusOK, scrapbooksResponse{Scrapbooks: scrapbooksResp})
}

// handlePostScrapOnScrapbook godoc
//
//	@summary post scrap on scrapbook
//	@description post scrap on scrapbook
//	@tags scraps
//	@security AccessTokenAuth
//	@param scrapID path string true "scrapID"
//	@param scrapbookID path string true "scrapbookID"
//	@success 201
//	@failure 400 {object} errorResponse
//	@failure 500 {object} errorResponse
//	@router /me/scraps/{scrapID}/scrapbooks/{scrapbookID} [post]
func (s *Server) handlePostScrapOnScrapbook(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	scrapID := ctx.Param("scrapID")
	scrapbookID := ctx.Param("scrapbookID")

	if err := s.db.InsertScrapOnScrapbook(ctx, userID, scrapID, scrapbookID); err != nil {
		golog.Error("handlePostMyScrapOnScrapbook: insert scrap on scrapbook: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusCreated)
}

// handleDeleteScrapOnScrapbook godoc
//
//	@summary delete scrap on scrapbook
//	@description delete scrap on scrapbook
//	@tags scraps
//	@security AccessTokenAuth
//	@param scrapID path string true "scrapID"
//	@param scrapbookID path string true "scrapbookID"
//	@success 204
//	@failure 400 {object} errorResponse
//	@failure 500 {object} errorResponse
//	@router /me/scraps/{scrapID}/scrapbooks/{scrapbookID} [delete]
func (s *Server) handleDeleteScrapOnScrapbook(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	scrapID := ctx.Param("scrapID")
	scrapbookID := ctx.Param("scrapbookID")

	if err := s.db.DeleteScrapOnScrapbook(ctx, userID, scrapID, scrapbookID); err != nil {
		golog.Error("handleDeleteScrapOnScrapbook: delete scrap on scrapbook: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
