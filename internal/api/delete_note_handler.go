package api

import (
	"context"
	response "github.com/Khasmag06/kode-notes/pkg/http"
	"github.com/go-chi/chi/v5"
	"net/http"
)

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
