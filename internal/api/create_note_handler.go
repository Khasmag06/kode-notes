package api

import (
	"context"
	"encoding/json"
	"github.com/Khasmag06/kode-notes/internal/models"
	response "github.com/Khasmag06/kode-notes/pkg/http"
	"net/http"
)

// CreateNote
// @Tags Note
// @Summary create a new note
// @Description create a new note
// @ID createNote
// @Accept  json
// @Produce json
// @Param input body models.Note true "note info"
// @Success 200 {object} http.SuccessResponse
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Security ApiKeyAuth
// @Router /note/create [post]
func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIdParam).(int)
	var note models.Note

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}
	defer r.Body.Close()

	if err := h.validator.Struct(note); err != nil {
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
