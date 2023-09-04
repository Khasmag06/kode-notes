package api

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"

	response "github.com/Khasmag06/kode-notes/pkg/http"
)

// GetNote
// @Tags Note
// @Summary get note
// @Description get note
// @ID getNote
// @Accept  json
// @Produce json
// @Param note_id path int64 true "ID of the note to retrieve"
// @Success 200 {object} http.SuccessResponse
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Security ApiKeyAuth
// @Router /note/get/{note_id} [get]
func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIdParam).(int)

	ctx := context.Background()

	noteId, err := parseNoteId(chi.URLParam(r, "note_id"))
	if err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}

	note, err := h.noteService.GetNote(ctx, userID, noteId)
	if err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}

	response.WriteSuccessResponse(w, note)
}
