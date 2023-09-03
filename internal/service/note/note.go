package note

import (
	"context"
	"github.com/Khasmag06/kode-notes/internal/models"

	businessErr "github.com/Khasmag06/kode-notes/pkg/app_err"
)

type repository interface {
	CreateNote(ctx context.Context, userID int, note models.Note) error
	UpdateNote(ctx context.Context, userID int, note models.Note) error
	DeleteNote(ctx context.Context, userID int, noteID int) error
	GetAllNotes(ctx context.Context, userID int) ([]models.Note, error)
	GetNote(ctx context.Context, userID int, noteID int) (*models.Note, error)
}

type service struct {
	repo repository
}

const (
	noteNotExistsErr = "Заметка не найдена"
)

func New(r repository) *service {
	return &service{
		repo: r,
	}
}

func (s *service) CreateNote(ctx context.Context, userID int, note models.Note) error {
	return s.repo.CreateNote(ctx, userID, note)
}

func (s *service) UpdateNote(ctx context.Context, userID int, data models.Note) error {
	note, err := s.repo.GetNote(ctx, userID, data.ID)
	if err != nil {
		return err
	}

	if note == nil {
		return businessErr.NewBusinessError(noteNotExistsErr)
	}

	return s.repo.UpdateNote(ctx, userID, data)
}

func (s *service) DeleteNote(ctx context.Context, userID int, noteID int) error {
	note, err := s.repo.GetNote(ctx, userID, noteID)
	if err != nil {
		return err
	}

	if note == nil {
		return businessErr.NewBusinessError(noteNotExistsErr)
	}

	return s.repo.DeleteNote(ctx, userID, noteID)
}

func (s *service) GetAllNotes(ctx context.Context, userID int) ([]models.Note, error) {
	notes, err := s.repo.GetAllNotes(ctx, userID)
	if err != nil {
		return nil, err
	}
	if notes == nil {
		return []models.Note{}, nil
	}
	return notes, nil
}

func (s *service) GetNote(ctx context.Context, userID int, noteID int) (models.Note, error) {
	note, err := s.repo.GetNote(ctx, userID, noteID)
	if err != nil {
		return models.Note{}, err
	}
	if note == nil {
		return models.Note{}, businessErr.NewBusinessError(noteNotExistsErr)
	}

	return *note, nil
}
