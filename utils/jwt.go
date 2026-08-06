package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID       int `json:"userId"`
	TokenVersion int `json:"tokenVersion"`
	TokenID      int `json:"tokenId"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, tokenVersion, tokenID int) (string, error) {

	claims := JWTClaims{
		UserID:       userID,
		TokenVersion: tokenVersion,
		TokenID:      tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}