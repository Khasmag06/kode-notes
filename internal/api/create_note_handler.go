package api

import (
	"context"
	"encoding/json"
	"github.com/Khasmag06/kode-notes/internal/models"
	businessErr "github.com/Khasmag06/kode-notes/pkg/app_err"
	response "github.com/Khasmag06/kode-notes/pkg/http"
	"net/http"
)

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIdParam).(int)
	var note models.Note

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}
	defer r.Body.Close()

	if err := checkCreateRequestData(note); err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}

	if err := h.CheckSpellingNote(note); err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}

	ctx := context.Background()

	if err := h.noteService.CreateNote(ctx, userID, note); err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}

	response.WriteSuccessResponse(w, nil)
}

func checkCreateRequestData(note models.Note) error {
	if note.Title == "" && note.Content == "" {
		return businessErr.NewBusinessError(noteIsEmptyErr)
	}

	return nil
}
