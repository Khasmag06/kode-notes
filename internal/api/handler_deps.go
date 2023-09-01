package api

import (
	"context"
	"github.com/Khasmag06/kode-notes/internal/models"
)

type authService interface {
	SignUp(ctx context.Context, user models.User) error
	Login(ctx context.Context, loginData models.User) (*models.TokensResponse, error)

	GenerateToken(userHash string) (string, error)
	ParseToken(accessToken string) (*models.TokenClaims, error)
}

type NoteService interface {
	CreateNote(ctx context.Context, userID int, note models.Note) error
	UpdateNote(ctx context.Context, userID int, note models.Note) error
	DeleteNote(ctx context.Context, userID int, noteId int) error
	GetAllNotes(ctx context.Context, userID int) ([]models.Note, error)
	GetNote(ctx context.Context, userID int, noteId int) (models.Note, error)
}

type Logger interface {
	Info(text ...any)
	Warn(text ...any)
	Error(text ...any)
}

type decoder interface {
	Encrypt(data []byte) (string, error)
	Decrypt(encrypted string) ([]byte, error)
}

type speller interface {
	CheckText(text string) ([]models.SpellError, error)
}
