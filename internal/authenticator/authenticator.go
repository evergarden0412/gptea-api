package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type Authenticator struct {
}

type AccessToken struct {
	jwt.RegisteredClaims
	UserID string `json:"userID"`
}

type RefreshToken struct {
	jwt.RegisteredClaims
	AccessTokenID string `json:"ati"`
}
