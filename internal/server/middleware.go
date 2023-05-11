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
		golog.Error("ensureUser: bind header:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	atStr, found := strings.CutPrefix(header.Authorization, "Bearer ")
	if !found {
		golog.Error("ensureUser: cut prefix: not found")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{Error: "no bearer prefix"})
		return
	}
	if atStr == "" {
		golog.Error("ensureUser: no access token")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{Error: "no access token"})
		return
	}

	at, err := s.a.VerifyAccessToken(atStr)
	if err != nil {
		golog.Error("ensureUser: verify access token:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	ctx.Set(ContextKeyUserID, at.Subject)
}
