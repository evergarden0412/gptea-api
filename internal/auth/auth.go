package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidSignMethod = errors.New("invalid signing method")
	ErrTokensNotMatch    = errors.New("tokens not match")
)

type AccessToken struct {
	jwt.RegisteredClaims
	signed string
}

func (a AccessToken) Signed() string {
	return a.signed
}

type RefreshToken struct {
	jwt.RegisteredClaims
	AccessTokenID string `json:"ati"`
	signed        string
}

func (r RefreshToken) Signed() string {
	return r.signed
}

type AuthenticatorConfig struct {
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	AccessTokenKey  []byte
	RefreshTokenKey []byte
}

type Authenticator struct {
	cfg AuthenticatorConfig
}

func New(cfg AuthenticatorConfig) *Authenticator {
	return &Authenticator{cfg: cfg}
}

func (a *Authenticator) accessTokenKeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, ErrInvalidSignMethod
	}
	return a.cfg.AccessTokenKey, nil
}

func (a *Authenticator) refreshTokenKeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, ErrInvalidSignMethod
	}
	return a.cfg.RefreshTokenKey, nil
}

func (a *Authenticator) IssueAccessToken(userID string) (AccessToken, error) {
	id := make([]byte, 15) // base64 multiple of 3
	if _, err := rand.Read(id); err != nil {
		return AccessToken{}, err
	}

	at := AccessToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(a.cfg.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Subject:   userID,
			ID:        base64.RawStdEncoding.EncodeToString(id),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, at)
	signed, err := token.SignedString(a.cfg.AccessTokenKey)
	if err != nil {
		return AccessToken{}, err
	}
	at.signed = signed
	return at, nil
}

func (a *Authenticator) IssueRefreshToken(accessTokenID string) (RefreshToken, error) {
	id := make([]byte, 15) // base64 multiple of 3
	if _, err := rand.Read(id); err != nil {
		return RefreshToken{}, err
	}

	rt := RefreshToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(a.cfg.RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ID:        base64.RawStdEncoding.EncodeToString(id),
		},
		AccessTokenID: accessTokenID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, rt)
	signed, err := token.SignedString(a.cfg.RefreshTokenKey)
	if err != nil {
		return RefreshToken{}, err
	}
	rt.signed = signed
	return rt, nil
}

func (a *Authenticator) VerifyAccessToken(token string) (AccessToken, error) {
	var at AccessToken
	if _, err := jwt.ParseWithClaims(token, &at, a.accessTokenKeyFunc); err != nil {
		return AccessToken{}, err
	}
	return at, nil
}

func (a *Authenticator) VerifyAccessTokenForRefresh(token string) (AccessToken, error) {
	var at AccessToken
	_, err := jwt.ParseWithClaims(token, &at, a.accessTokenKeyFunc)
	if err != nil && errors.Is(err, jwt.ErrTokenExpired) {
		return AccessToken{}, err
	}
	return at, nil
}

func (a *Authenticator) VerifyRefreshToken(token string) (RefreshToken, error) {
	var rt RefreshToken
	if _, err := jwt.ParseWithClaims(token, &rt, a.refreshTokenKeyFunc); err != nil {
		return RefreshToken{}, err
	}
	return rt, nil
}

func (a *Authenticator) RefreshAccessToken(accessToken AccessToken, refreshToken RefreshToken) (AccessToken, RefreshToken, error) {
	if accessToken.ID != refreshToken.AccessTokenID {
		return AccessToken{}, RefreshToken{}, ErrTokensNotMatch
	}

	newAT, err := a.IssueAccessToken(accessToken.Subject)
	if err != nil {
		return AccessToken{}, RefreshToken{}, err
	}
	newRT, err := a.IssueRefreshToken(newAT.ID)
	if err != nil {
		return AccessToken{}, RefreshToken{}, err
	}

	return newAT, newRT, nil
}
