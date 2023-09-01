package api

import (
	"context"
	"github.com/Khasmag06/kode-notes/pkg/app_err"
	response "github.com/Khasmag06/kode-notes/pkg/http"
	"strconv"
	"strings"

	"net/http"
)

const (
	authorizationHeader = "Authorization"
	userIdParam         = "userId"
)

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get(authorizationHeader), "Bearer ")
		if len(authHeader) != 2 {
			response.WriteErrorResponse(w, h.logger, app_err.NewUnauthorizedError())
			h.logger.Error(`header 'Authorization' невалиден`)
			return
		}

		accessToken := authHeader[1]

		claims, err := h.authService.ParseToken(accessToken)
		if err != nil {
			response.WriteErrorResponse(w, h.logger, app_err.NewUnauthorizedError())
			h.logger.Error(err)
			return
		}

		userIdBytes, err := h.decoder.Decrypt(claims.UserID)
		if err != nil {
			response.WriteErrorResponse(w, h.logger, err)
			return
		}
		userID, _ := strconv.Atoi(string(userIdBytes))

		ctx := context.WithValue(r.Context(), userIdParam, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
