package api

import (
	"context"
	response "github.com/Khasmag06/kode-notes/pkg/http"
	"net/http"
)

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
