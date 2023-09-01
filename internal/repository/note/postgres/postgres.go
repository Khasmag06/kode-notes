package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Khasmag06/kode-notes/internal/models"
)

type repo struct {
	pool *pgxpool.Pool
}

func New(db *pgxpool.Pool) *repo {
	return &repo{
		pool: db,
	}
}

func (r *repo) CreateNote(ctx context.Context, userID int, note models.Note) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO notes (user_id, title, content)
			VALUES ($1, $2, $3)`, userID, note.Title, note.Content)
	if err != nil {
		return fmt.Errorf("SQL: create note: %w", err)
	}

	return nil
}

func (r *repo) UpdateNote(ctx context.Context, userID int, note models.Note) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE notes 
			SET title = $1, content = $2
			WHERE id = $3 AND user_id = $4`, note.Title, note.Content, note.ID, userID)
	if err != nil {
		return fmt.Errorf("SQL: update note: %w", err)
	}

	return nil
}

func (r *repo) DeleteNote(ctx context.Context, userID int, noteID int) error {
	_, err := r.pool.Exec(ctx,
		`DELETE 
			FROM notes
			WHERE id = $1 AND user_id = $2`, noteID, userID)
	if err != nil {
		return fmt.Errorf("SQL: delete note: %w", err)
	}

	return nil
}

func (r *repo) GetAllNotes(ctx context.Context, userID int) ([]models.Note, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT 
			id, title, content
			FROM notes WHERE user_id = $1`, userID)

	if err != nil {
		return nil, fmt.Errorf("SQL: get all notes: %w", err)
	}
	defer rows.Close()

	var notes []models.Note

	for rows.Next() {
		var note models.Note

		err := rows.Scan(&note.ID, &note.Title, &note.Content)
		if err != nil {
			return nil, fmt.Errorf("SQL: scan row in get all notes: %w", err)
		}

		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("SQL: iterate rows in get all notes: %w", err)
	}

	return notes, nil
}

func (r *repo) GetNote(ctx context.Context, userID int, noteID int) (*models.Note, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT 
			id, title, content
			FROM notes
			WHERE id = $1 AND user_id = $2`, noteID, userID)

	var note models.Note

	err := row.Scan(&note.ID, &note.Title, &note.Content)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("SQL: get note: %w", err)
	}

	return &note, nil
}
