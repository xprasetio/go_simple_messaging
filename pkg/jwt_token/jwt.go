package jwttoken

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kooroshh/fiber-boostrap/pkg/env"
)

type ClaimToken struct {
	Username string `json:"username"`
	Fullname string `json:"full_name"`
	jwt.RegisteredClaims
}

var mapTypeToken = map[string]time.Duration{
	"token":         time.Hour * 3,      // 3 hours
	"refresh_token": time.Hour * 24 * 7, // 7 days
}

func GenerateToken(ctx context.Context, username string, fullname string, tokenType string) (string, error) {
	secret := []byte(env.GetEnv("APP_JWT_SECRET", ""))

	claimToken := ClaimToken{
		Username: username,
		Fullname: fullname,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    env.GetEnv("APP_NAME", ""),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(mapTypeToken[tokenType])),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimToken)
	resultToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return resultToken, fmt.Errorf("failed to sign token: %w", err)
	}
	return resultToken, nil
}
