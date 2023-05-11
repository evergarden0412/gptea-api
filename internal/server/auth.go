package server

import (
	"net/http"
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
	userID, err := internal.NewUserID()
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
	if err != nil {
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
	ctx.JSON(http.StatusOK, signInHandlerOutput{
		AccessToken:  at.Signed(),
		RefreshToken: rt.Signed(),
	})
}
