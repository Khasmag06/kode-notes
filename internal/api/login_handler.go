package api

import (
	"context"
	"encoding/json"
	"github.com/Khasmag06/kode-notes/internal/models"
	response "github.com/Khasmag06/kode-notes/pkg/http"
	"net/http"
	"strings"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq models.User
	ctx := context.Background()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginReq); err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}

	loginReq.Login = strings.TrimSpace(strings.ToLower(loginReq.Login))

	tokenData, err := h.authService.Login(ctx, loginReq)
	if err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}

	response.WriteSuccessResponse(w, tokenData)
}
