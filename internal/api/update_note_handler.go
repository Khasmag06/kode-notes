package api

import (
	"context"
	"encoding/json"
	"github.com/Khasmag06/kode-notes/internal/models"
	"net/http"

	businessErr "github.com/Khasmag06/kode-notes/pkg/app_err"
	response "github.com/Khasmag06/kode-notes/pkg/http"
)

func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIdParam).(int)
	var note models.Note
	urlParam, err := parseNoteId(r.URL.Query().Get("noteId"))
	if err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}
	note.ID = urlParam

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&note); err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}
	defer r.Body.Close()

	err = checkUpdateRequestNote(note)
	if err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}

	if err := h.CheckSpellingNote(note); err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}

	ctx := context.Background()

	err = h.noteService.UpdateNote(ctx, userID, note)
	if err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}

	response.WriteSuccessResponse(w, nil)
}

func checkUpdateRequestNote(note models.Note) error {
	if note.ID == 0 {
		return businessErr.NewBusinessError(noteIsEmptyErr)
	}
	if note.Title == "" && note.Content == "" {
		return businessErr.NewBusinessError(noteIsEmptyErr)
	}

	return nil
}
