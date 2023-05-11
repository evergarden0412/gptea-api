package credential

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Credential interface {
	Verify(ctx context.Context, token string) (VerifyResult, error)
}

type VerifyBody struct {
	Token              string `json:"token" binding:"required"`
	CredentialProvider string `json:"credentialProvider" binding:"required" example:"naver"`
}

type VerifyResult struct {
	CredentialProvider string
	CredentialID       string
}

const (
	ProviderNaver = "naver"
	ProviderKakao = "kakao"
)

var (
	ErrUnknownProvider = fmt.Errorf("unknown provider")
)

func New(provider string) (Credential, error) {
	switch provider {
	case ProviderNaver:
		return &naverCredential{}, nil
	case ProviderKakao:
		return &kakaoCredential{}, nil
	default:
		return nil, ErrUnknownProvider
	}
}

type naverCredential struct {
}

type kakaoCredential struct {
}

type naverProfileResponse struct {
	Resultcode string `json:"resultcode"`
	Message    string `json:"message"`
	Response   struct {
		ID string `json:"id" binding:"required"`
	}
}

const naverProfileURL = "https://openapi.naver.com/v1/nid/me"

func (c *naverCredential) Verify(ctx context.Context, token string) (VerifyResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, naverProfileURL, nil)
	if err != nil {
		return VerifyResult{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return VerifyResult{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return VerifyResult{}, err
	}
	var naverProfile naverProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&naverProfile); err != nil {
		return VerifyResult{}, err
	}
	if naverProfile.Response.ID == "" {
		return VerifyResult{}, err
	}
	return VerifyResult{
		CredentialProvider: ProviderNaver,
		CredentialID:       naverProfile.Response.ID,
	}, nil
}

type kakaoProfileResponse struct {
	ID int64 `json:"id" binding:"required"`
}

const kakaoProfileURL = "https://kapi.kakao.com/v2/user/me"

func (c *kakaoCredential) Verify(ctx context.Context, token string) (VerifyResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, kakaoProfileURL, nil)
	if err != nil {
		return VerifyResult{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return VerifyResult{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return VerifyResult{}, err
	}
	var kakaoProfile kakaoProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&kakaoProfile); err != nil {
		return VerifyResult{}, err
	}
	if kakaoProfile.ID == 0 {
		return VerifyResult{}, err
	}
	return VerifyResult{
		CredentialProvider: ProviderKakao,
		CredentialID:       fmt.Sprint(kakaoProfile.ID),
	}, nil
}
