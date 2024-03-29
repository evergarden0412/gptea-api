package server

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/evergarden0412/gptea-api/internal"
	"github.com/evergarden0412/gptea-api/internal/credential"
	"github.com/evergarden0412/gptea-api/internal/postgres"
	"github.com/gin-gonic/gin"
	"github.com/kataras/golog"
)

type credBody struct {
	Cred        string `json:"cred" binding:"required" example:"naver"`
	AccessToken string `json:"accessToken" binding:"required" `
}

// handleRegister godoc
// @Summary Register a credential
// @Description Register a credential
// @Param body body credBody true "body"
// @Success 201 {object} messageResponse
// @Router /auth/cred/register [post]
// @Tags auth
func (s *Server) handleRegister(ctx *gin.Context) {
	var body credBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		golog.Error("handleRegister: bind json: ", err)
		return
	}

	cred, err := credential.New(body.Cred)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		golog.Error("handleRegister: new credential: ", err)
		return
	}
	verifyResult, err := cred.Verify(ctx, body.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		golog.Error("handleRegister: verify: ", err)
		return
	}
	userID, err := internal.NewID()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		golog.Error("handleRegister: new user id: ", err)
		return
	}
	now := time.Now().UTC()
	if err := s.db.Register(ctx, postgres.RegisterInput{
		UserID:         userID,
		CredentialType: verifyResult.CredentialProvider,
		CredentialID:   verifyResult.CredentialID,
		CreatedAt:      &now,
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		golog.Error("handleRegister: register: ", err)
		return
	}
	ctx.JSON(http.StatusCreated, messageResponse{Message: "success"})
}

type signInHandlerOutput struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// handleSignIn godoc
// @Summary Sign in with a credential
// @Description Sign in with a credential
// @Param body body credBody true "body"
// @Success 200 {object} signInHandlerOutput
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /auth/cred/sign-in [post]
// @Tags auth
func (s *Server) handleSignIn(ctx *gin.Context) {
	var body credBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	cred, err := credential.New(body.Cred)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		golog.Error("handleSignIn: new credential: ", err)
		return
	}
	verifyResult, err := cred.Verify(ctx, body.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		golog.Error("handleSignIn: verify: ", err)
		return
	}

	userID, err := s.db.SignIn(ctx, verifyResult.CredentialProvider, verifyResult.CredentialID)
	if errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		golog.Error("handleSignIn: sign in: ", err)
		return
	}

	at, err := s.a.IssueAccessToken(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		golog.Error("handleSignIn: issue access token: ", err)
		return
	}
	rt, err := s.a.IssueRefreshToken(at.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		golog.Error("handleSignIn: issue refresh token: ", err)
		return
	}
	if err := s.db.UpsertRefreshToken(ctx, at.Subject, rt.ID); err != nil {
		golog.Error("handleSignIn: upsert refresh token: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, signInHandlerOutput{
		AccessToken:  at.Signed(),
		RefreshToken: rt.Signed(),
	})
}

type tokenResponse struct {
	ExpiresAt time.Time `json:"exp,omitempty"`
	IssuedAt  time.Time `json:"iat,omitempty"`
	ID        string    `json:"jti,omitempty"`
	Subject   string    `json:"sub,omitempty"`
}

// handleRefreshToken godoc
// @Summary Refresh a token
// @Description Refresh a token
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Success 200 {object} signInHandlerOutput
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /auth/token/refresh [post]
// @tags token
func (s *Server) handleRefreshToken(ctx *gin.Context) {
	var header tokenHeader
	if err := ctx.ShouldBindHeader(&header); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		golog.Error("handleRefreshToken: bind header: ", err)
		return
	}
	atStr, found := strings.CutPrefix(header.Authorization, "Bearer ")
	if !found {
		ctx.JSON(http.StatusBadRequest, errorResponse{Error: "no bearer prefix"})
		golog.Error("handleRefreshToken: cut prefix: not found")
		return
	}
	if atStr == "" || header.XRefreshToken == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse{Error: "no token"})
		golog.Error("handleRefreshToken: no token")
		return
	}

	at, err := s.a.VerifyAccessTokenForRefresh(atStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		golog.Error("handleRefreshToken: verify access token: ", err)
		return
	}
	rt, err := s.a.VerifyRefreshToken(header.XRefreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		golog.Error("handleRefreshToken: verify refresh token: ", err)
		return
	}
	if exists, err := s.db.IsRefreshTokenExists(ctx, at.Subject, rt.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		golog.Error("handleRefreshToken: exists refresh token: ", err)
		return
	} else if !exists {
		ctx.JSON(http.StatusBadRequest, errorResponse{Error: "refresh token not in db"})
		golog.Error("handleRefreshToken: refresh token not in db")
		return
	}

	newAT, newRT, err := s.a.RefreshAccessToken(at, rt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		golog.Error("handleRefreshToken: refresh: ", err)
		return
	}
	if err := s.db.UpsertRefreshToken(ctx, newAT.Subject, newRT.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		golog.Error("handleRefreshToken: upsert refresh token: ", err)
		return
	}

	ctx.JSON(http.StatusOK, signInHandlerOutput{
		AccessToken:  newAT.Signed(),
		RefreshToken: newRT.Signed(),
	})
}

// handleLogout godoc
// @summary Logout
// @description delete refresh token
// @tags token
// @Security AccessTokenAuth
// @success 204
// @failure 400 {object} errorResponse
// @failure 500 {object} errorResponse
// @router /auth/cred/logout [delete]
func (s *Server) handleLogout(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	if err := s.db.Logout(ctx, userID); err != nil {
		golog.Error("handleLogout: delete refresh token: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// handleDeleteMe godoc
// @Summary Delete my account
// @Description Delete my account
// @tags users
// @Security AccessTokenAuth
// @Success 204
// @Failure 500 {object} errorResponse
// @Router /me [delete]
func (s *Server) handleDeleteMe(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	if err := s.db.Resign(ctx, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		golog.Error("handleDeleteMe: delete user: ", err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
