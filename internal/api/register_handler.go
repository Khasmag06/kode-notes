package api

import (
	"context"
	"encoding/json"
	"github.com/Khasmag06/kode-notes/internal/models"
	response "github.com/Khasmag06/kode-notes/pkg/http"
	"net/http"
	"strings"
)

// SignUp
// @Tags Auth
// @Summary SignUp
// @Description login
// @ID login
// @Accept  json
// @Produce json
// @Param input body models.User true "account info"
// @Success 200 {object} http.SuccessResponse
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /auth/login [post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var signUpReq models.User
	ctx := context.Background()

	if err := json.NewDecoder(r.Body).Decode(&signUpReq); err != nil {
		response.WriteErrorResponse(w, h.logger, err)
		return
	}

	if err := h.validator.Struct(signUpReq); err != nil {
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
