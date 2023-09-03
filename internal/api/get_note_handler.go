package api

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"

	response "github.com/Khasmag06/kode-notes/pkg/http"
)

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
