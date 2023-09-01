package models

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type User struct {
	Id        int       `json:"-"`
	Login     string    `json:"login" binding:"required"`
	Password  string    `json:"password" binding:"required,password,min=8,max=32"`
	CreatedAt time.Time `json:"-"`
}

type TokenClaims struct {
	jwt.RegisteredClaims
	UserID string
}

type TokensResponse struct {
	AccessToken string `json:"accessToken"`
}