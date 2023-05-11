package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kataras/golog"
)

const (
	ContextKeyUserID = "userID"
)

type tokenHeader struct {
	Authorization string `header:"authorization"`
	XRefreshToken string `header:"x-refresh-token"`
}

func (s *Server) ensureUser(ctx *gin.Context) {
	var header tokenHeader
	if err := ctx.ShouldBindHeader(&header); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		golog.Error("ensureUser: bind header:", err)
		ctx.Abort()
		return
	}
	atStr, found := strings.CutPrefix(header.Authorization, "Bearer ")
	if !found {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{Error: "no bearer prefix"})
		golog.Error("ensureUser: cut prefix: not found")
		ctx.Abort()
		return
	}
	if atStr == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{Error: "no access token"})
		golog.Error("ensureUser: no access token")
		ctx.Abort()
		return
	}

	at, err := s.a.VerifyAccessToken(atStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{Error: err.Error()})
		golog.Error("ensureUser: verify access token:", err)
		ctx.Abort()
		return
	}

	ctx.Set(ContextKeyUserID, at.Subject)
}
