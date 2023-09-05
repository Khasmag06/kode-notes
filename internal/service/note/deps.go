//go:generate mockgen -source=$GOFILE -destination=mocks_test.go -package=$GOPACKAGE
package note

import (
	"context"
	"github.com/Khasmag06/kode-notes/internal/models"
)

type repository interface {
	CreateNote(ctx context.Context, userID int, note models.Note) error
	UpdateNote(ctx context.Context, userID int, note models.Note) error
	DeleteNote(ctx context.Context, userID int, noteID int) error
	GetAllNotes(ctx context.Context, userID int) ([]models.Note, error)
	GetNote(ctx context.Context, userID int, noteID int) (*models.Note, error)
}
