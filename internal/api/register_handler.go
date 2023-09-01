package api

import (
	"context"
	"encoding/json"
	"github.com/Khasmag06/kode-notes/internal/models"
	response "github.com/Khasmag06/kode-notes/pkg/http"
	"net/http"
	"strings"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var signUpReq models.User
	ctx := context.Background()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&signUpReq); err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}

	signUpReq.Login = strings.TrimSpace(strings.ToLower(signUpReq.Login))

	if err := h.authService.SignUp(ctx, signUpReq); err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}

	response.WriteSuccessResponse(w, nil)
}
