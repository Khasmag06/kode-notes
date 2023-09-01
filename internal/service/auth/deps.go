package auth

import (
	"context"
	"github.com/Khasmag06/kode-notes/internal/models"
)

type repository interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
}

type hasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type jwt interface {
	GenerateToken(userHash string) (string, error)
	ParseToken(accessToken string) (*models.TokenClaims, error)
}

type decoder interface {
	Encrypt(data []byte) (string, error)
}
