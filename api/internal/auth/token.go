package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"chaseapp.tv/api/internal/config"
)

// ChatTokenSigner creates signed JWTs for chat authentication.
type ChatTokenSigner struct {
	cfg config.ChatConfig
}

// NewChatTokenSigner creates a new signer.
func NewChatTokenSigner(cfg config.ChatConfig) (*ChatTokenSigner, error) {
	if cfg.SigningKey == "" {
		return nil, errors.New("chat signing key is required")
	}
	if cfg.TokenTTL <= 0 {
		cfg.TokenTTL = 15 * time.Minute
	}
	return &ChatTokenSigner{cfg: cfg}, nil
}

// ChatClaims are the claims embedded in the chat token.
type ChatClaims struct {
	Issuer      string    `json:"iss,omitempty"`
	Audience    string    `json:"aud,omitempty"`
	Subject     string    `json:"sub"`
	Email       string    `json:"email,omitempty"`
	Permissions []string  `json:"permissions"`
	IssuedAt    time.Time `json:"iat"`
	ExpiresAt   time.Time `json:"exp"`
}

// Sign builds an HS256 JWT for the provided user and permissions.
func (s *ChatTokenSigner) Sign(userID, email string, permissions []string) (string, error) {
	if userID == "" {
		return "", errors.New("user id is required")
	}
	if len(permissions) == 0 {
		return "", errors.New("at least one permission is required")
	}

	now := time.Now().UTC()
	claims := ChatClaims{
		Issuer:      s.cfg.Issuer,
		Audience:    s.cfg.Audience,
		Subject:     userID,
		Email:       email,
		Permissions: permissions,
		IssuedAt:    now,
		ExpiresAt:   now.Add(s.cfg.TokenTTL),
	}

	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("marshal header: %w", err)
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("marshal claims: %w", err)
	}

	headerEnc := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsEnc := base64.RawURLEncoding.EncodeToString(claimsJSON)
	unsigned := headerEnc + "." + claimsEnc

	mac := hmac.New(sha256.New, []byte(s.cfg.SigningKey))
	mac.Write([]byte(unsigned))
	signature := mac.Sum(nil)
	sigEnc := base64.RawURLEncoding.EncodeToString(signature)

	return unsigned + "." + sigEnc, nil
}
