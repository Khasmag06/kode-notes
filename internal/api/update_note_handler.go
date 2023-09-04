package api

import (
	"context"
	"encoding/json"
	"github.com/Khasmag06/kode-notes/internal/models"
	"github.com/go-chi/chi/v5"
	"net/http"

	response "github.com/Khasmag06/kode-notes/pkg/http"
)

// UpdateNote
// @Tags Note
// @Summary update note
// @Description update note
// @ID updateNote
// @Accept  json
// @Produce json
// @Param note_id path int64 true "ID of the note to update"
// @Success 200 {object} http.SuccessResponse
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Security ApiKeyAuth
// @Router /note/update/{note_id} [put]
func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIdParam).(int)
	var note models.Note
	noteId, err := parseNoteId(chi.URLParam(r, "note_id"))
	note.ID = noteId

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&note); err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}
	defer r.Body.Close()

	err = checkRequestData(note)
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
