package api

import (
	"context"
	response "github.com/Khasmag06/kode-notes/pkg/http"
	"net/http"
)

// GetAllNotes
// @Tags Note
// @Summary get all notes
// @Description get all notes
// @ID getAllNotes
// @Accept  json
// @Produce json
// @Success 200 {object} http.SuccessResponse
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Security ApiKeyAuth
// @Router /note/get-all [get]
func (h *Handler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIdParam).(int)

	ctx := context.Background()

	notes, err := h.noteService.GetAllNotes(ctx, userID)
	if err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}

	response.WriteSuccessResponse(w, notes)

}
