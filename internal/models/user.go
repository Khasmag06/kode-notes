package models

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type User struct {
	Id        int       `json:"-"`
	Username  string    `json:"username" validate:"required" format:"string" example:"my_login"`
	Password  string    `json:"password" validate:"required,password,min=8,max=32" format:"string" example:"Qwerty123!"`
	CreatedAt time.Time `json:"-"`
}

type TokenClaims struct {
	jwt.RegisteredClaims
	UserID string
}

type TokensResponse struct {
	AccessToken string `json:"accessToken"`
}
