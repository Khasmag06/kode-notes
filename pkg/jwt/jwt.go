package jwt

import (
	"errors"
	"fmt"
	"github.com/Khasmag06/kode-notes/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	ErrEmptySigningKey = errors.New("empty signing key")
	ErrTokenSigning    = errors.New("failed to sign token")
	ErrParseToken      = errors.New("failed to parse token")
	ErrInvalidToken    = errors.New("invalid token claims")
)

type JWT struct {
	signKey string
}

func New(signKey string) (*JWT, error) {
	if signKey == "" {
		return nil, ErrEmptySigningKey
	}

	return &JWT{signKey: signKey}, nil
}

func (j *JWT) GenerateToken(userHash string) (string, error) {
	// generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userHash,
	})

	// sign token
	tokenString, err := token.SignedString([]byte(j.signKey))
	if err != nil {
		return "", ErrTokenSigning
	}
	return tokenString, nil
}

func (j *JWT) ParseToken(accessToken string) (*models.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.signKey), nil
	})

	if err != nil {
		return nil, ErrParseToken
	}

	claims, ok := token.Claims.(*models.TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
