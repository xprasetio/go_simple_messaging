package jwt_token

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

var MapTypeToken = map[string]time.Duration{
	"token":         time.Hour * 3,      // 3 hours
	"refresh_token": time.Hour * 24 * 3, // 3 days
}

var jwtSecret = []byte(env.GetEnv("APP_SECRET", ""))

func GenerateToken(ctx context.Context, username string, fullname string, tokenType string, now time.Time) (string, error) {
	// secret := []byte(env.GetEnv("APP_SECRET", ""))

	claimToken := ClaimToken{
		Username: username,
		Fullname: fullname,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    env.GetEnv("APP_NAME", ""),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(MapTypeToken[tokenType])),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimToken)
	resultToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return resultToken, fmt.Errorf("failed to sign token: %V", err)
	}
	return resultToken, nil
}

func ValidateToken(ctx context.Context, token string) (*ClaimToken, error) {
	var (
		claimToken *ClaimToken
		ok         bool
	)
	jwtToken, err := jwt.ParseWithClaims(token, &ClaimToken{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}
	if claimToken, ok = jwtToken.Claims.(*ClaimToken); !ok || !jwtToken.Valid {
		return nil, fmt.Errorf("invalid tokenss")
	}
	return claimToken, nil
}
