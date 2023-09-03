package api

import (
	"fmt"
	"github.com/Khasmag06/kode-notes/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strconv"
	"strings"

	businessErr "github.com/Khasmag06/kode-notes/pkg/app_err"
)

type Handler struct {
	*chi.Mux
	authService authService
	noteService NoteService
	decoder     decoder
	logger      Logger
	speller     speller
}

const (
	noteIsEmptyErr   = "Заметка не может быть пустой!"
	invalidNoteIdErr = "Невалидный параметр noteId"
)

func NewHandler(auth authService, noteService NoteService, decoder decoder, logger Logger, speller speller) http.Handler {
	h := Handler{
		Mux:         chi.NewMux(),
		authService: auth,
		noteService: noteService,
		decoder:     decoder,
		logger:      logger,
		speller:     speller,
	}

	h.Use(middleware.Recoverer)
	h.Use(middleware.Logger)

	h.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-up", h.SignUp)
			r.Post("/login", h.Login)
		})
		r.Route("/note", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(h.authMiddleware)
				r.Post("/create", h.CreateNote)
				r.Put("/update/{note_id}", h.UpdateNote)
				r.Delete("/delete/{note_id}", h.DeleteNote)
				r.Get("/get-all", h.GetAllNotes)
				r.Get("/get/{note_id}", h.GetNote)
			})
		})
	})

	return h
}

func (h *Handler) CheckSpellingNote(note models.Note) error {
	text := fmt.Sprintf("%s %s", note.Title, note.Content)
	spellErrors, err := h.speller.CheckText(text)
	if err != nil {
		return fmt.Errorf("unable to check spelling: %w", err)
	}

	if len(spellErrors) > 0 {
		var errorMsg strings.Builder
		for _, errorItem := range spellErrors {
			suggestions := strings.Join(errorItem.S, ", ")
			errorStr := fmt.Sprintf("Ошибка в слове: '%s'\nВозможно мы имели в виду: '%s'\n\n", errorItem.Word, suggestions)
			errorMsg.WriteString(errorStr)
		}
		return businessErr.NewBusinessError(errorMsg.String())
	}
	return nil
}

func parseNoteId(noteIdQuery string) (int, error) {
	noteId, err := strconv.Atoi(noteIdQuery)
	if err != nil {
		return 0, businessErr.NewBusinessError(invalidNoteIdErr)
	}
	if noteId == 0 {
		return 0, businessErr.NewBusinessError(invalidNoteIdErr)
	}

	return noteId, nil
}

func checkRequestData(note models.Note) error {
	if note.Title == "" && note.Content == "" {
		return businessErr.NewBusinessError(noteIsEmptyErr)
	}

	return nil
}
