package server

import (
	"net/http"

	"github.com/evergarden0412/gptea-api/internal"
	"github.com/evergarden0412/gptea-api/internal/postgres"
	"github.com/gin-gonic/gin"
	"github.com/kataras/golog"
)

type chatBody struct {
	Name string `json:"name"`
}

// handlePostMyChat godoc
// @summary Post my chat
// @description Post my chat
// @tags chats
// @security AccessTokenAuth
// @param body body chatBody true "body"
// @success 201 {object} messageResponse
// @failure 400 {object} errorResponse
// @failure 500 {object} errorResponse
// @router /me/chats [post]
func (s *Server) handlePostMyChat(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	var body chatBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		golog.Error("handlePostMyChat: bind json: ", err)
		return
	}

	chat, err := internal.NewChat()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		golog.Error("handlePostMyChat: new chat: ", err)
		return
	}
	chat.Name = body.Name

	if err := s.db.InsertChat(ctx, userID, *chat); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		golog.Error("handlePostMyChat: insert chat: ", err)
		return
	}

	ctx.JSON(http.StatusCreated, messageResponse{Message: "ok"})
}

type chatsResponse struct {
	Chats []internal.Chat `json:"chats"`
}

// handleGetMyChats godoc
// @Summary Get my chats
// @Description Get my chats in descending order of created_at
// @Security AccessTokenAuth
// @Success 200 {object} chatsResponse
// @Failure 500 {object} errorResponse
// @Router /me/chats [get]
// @tags chats
func (s *Server) handleGetMyChats(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	chats, err := s.db.SelectMyChats(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		golog.Error("handleGetMyChats: select chats: ", err)
		return
	}
	chatsForResp := make([]internal.Chat, len(chats))
	copy(chatsForResp, chats)

	ctx.JSON(http.StatusOK, chatsResponse{Chats: chatsForResp})
}

// handlePatchMyChat godoc
// @summary Patch my chat
// @description Patch my chat name
// @tags chats
// @security AccessTokenAuth
// @param chatID path string true "chatID"
// @param body body chatBody true "body"
// @success 204
// @failure 400 {object} errorResponse
// @failure 401 {object} errorResponse
// @failure 500 {object} errorResponse
// @router /me/chats/{chatID} [patch]
func (s *Server) handlePatchMyChat(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	chatID := ctx.Param("chatID")
	var body chatBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		golog.Error("handlePatchMyChat: bind json: ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	if err := s.db.PatchChat(ctx, userID, internal.Chat{
		ID:   chatID,
		Name: body.Name,
	}); err != nil {
		golog.Error("handlePatchMyChat: update chat: ", err)
		switch err {
		case postgres.ErrUnauthorized:
			ctx.JSON(http.StatusUnauthorized, errorResponse{Error: err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

// handleDeleteMyChat godoc
// @summary Delete my chat
// @description Delete my chat
// @tags chats
// @security AccessTokenAuth
// @param chatID path string true "chatID"
// @success 204
// @failure 400 {object} errorResponse
// @failure 401 {object} errorResponse
// @failure 500 {object} errorResponse
// @router /me/chats/{chatID} [delete]
func (s *Server) handleDeleteMyChat(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	chatID := ctx.Param("chatID")

	if err := s.db.DeleteChat(ctx, userID, chatID); err != nil {
		golog.Error("handleDeleteMyChat: delete chat: ", err)
		switch err {
		case postgres.ErrUnauthorized:
			ctx.JSON(http.StatusUnauthorized, errorResponse{Error: err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}
