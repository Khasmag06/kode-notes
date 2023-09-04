package api

import (
	"context"
	response "github.com/Khasmag06/kode-notes/pkg/http"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// DeleteNote
// @Tags Note
// @Summary delete a note
// @Description delete note
// @ID deleteNote
// @Accept  json
// @Produce json
// @Param note_id path int64 true "ID of the note to delete"
// @Success 200 {object} http.SuccessResponse
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Security ApiKeyAuth
// @Router /note/delete/{note_id} [delete]
func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIdParam).(int)

	noteId, err := parseNoteId(chi.URLParam(r, "note_id"))
	if err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}
	ctx := context.Background()

	err = h.noteService.DeleteNote(ctx, userID, noteId)
	if err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}

	response.WriteSuccessResponse(w, nil)
}
