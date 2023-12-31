package api

import (
	"context"
	"encoding/json"
	"github.com/Khasmag06/kode-notes/internal/models"
	response "github.com/Khasmag06/kode-notes/pkg/http"
	"net/http"
	"strings"
)

// Login
// @Tags Auth
// @Summary Login
// @Description Create account
// @ID create-account
// @Accept  json
// @Produce json
// @Param input body models.User true "account info"
// @Success 200 {object} http.SuccessResponse
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /auth/sign-up [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq models.User
	ctx := context.Background()

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}

	if err := h.validator.Struct(loginReq); err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}

	loginReq.Username = strings.TrimSpace(strings.ToLower(loginReq.Username))

	tokenData, err := h.authService.Login(ctx, loginReq)
	if err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}

	response.WriteSuccessResponse(w, tokenData)
}
