package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Khasmag06/kode-notes/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type repo struct {
	repository
	redis  *redis.Client
	logger logger
}

func New(rdb *redis.Client, noteRepo repository, logger logger) *repo {
	return &repo{
		repository: noteRepo,
		redis:      rdb,
		logger:     logger,
	}
}

const expiration = 48 * time.Hour // two days

func (r *repo) GetAllNotes(ctx context.Context, userID int) ([]models.Note, error) {
	notesCache, err := r.GetNoteFromCache(ctx, userID)
	if err != nil {
		r.logger.Error(err)
	}
	if notesCache != nil {
		return notesCache, nil
	}

	notes, err := r.repository.GetAllNotes(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := r.SaveNoteToCache(ctx, userID, notes); err != nil {
		r.logger.Error(err)
	}
	return notes, nil
}

func (r *repo) CreateNote(ctx context.Context, userID int, note models.Note) error {
	if err := r.repository.CreateNote(ctx, userID, note); err != nil {
		return err
	}
	if err := r.DeleteNoteFromCache(ctx, userID); err != nil {
		r.logger.Error(err)
	}
	return nil
}

func (r *repo) UpdateNote(ctx context.Context, userID int, note models.Note) error {
	if err := r.repository.UpdateNote(ctx, userID, note); err != nil {
		return err
	}
	if err := r.DeleteNoteFromCache(ctx, userID); err != nil {
		r.logger.Error(err)
	}

	return nil
}

func (r *repo) DeleteNote(ctx context.Context, userID int, noteID int) error {
	if err := r.repository.DeleteNote(ctx, userID, noteID); err != nil {
		return err
	}
	if err := r.DeleteNoteFromCache(ctx, userID); err != nil {
		r.logger.Error(err)
	}
	return nil
}

func (r *repo) SaveNoteToCache(ctx context.Context, userID int, notes []models.Note) error {
	user := fmt.Sprintf("u:%d", userID) // u - user
	notesJSON, err := json.Marshal(notes)
	if err != nil {
		return err
	}
	if err := r.redis.Set(ctx, user, notesJSON, expiration).Err(); err != nil {
		return err
	}
	return nil
}

func (r *repo) GetNoteFromCache(ctx context.Context, userID int) ([]models.Note, error) {
	user := fmt.Sprintf("u:%d", userID) // u - user
	notesJSON, err := r.redis.Get(ctx, user).Result()
	if err != nil {
		return nil, err
	}

	var notesCache []models.Note
	if err := json.Unmarshal([]byte(notesJSON), &notesCache); err != nil {
		return nil, err
	}
	return notesCache, nil
}

func (r *repo) DeleteNoteFromCache(ctx context.Context, userID int) error {
	user := fmt.Sprintf("u:%d", userID)
	if err := r.redis.Del(ctx, user).Err(); err != nil {
		return err
	}
	return nil
}
