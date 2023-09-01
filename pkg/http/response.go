package http

import (
	"encoding/json"
	"errors"
	"github.com/Khasmag06/kode-notes/pkg/app_err"
	"github.com/Khasmag06/kode-notes/pkg/cust_validator"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type SuccessResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status string `json:"status"`
	Error  Error  `json:"error"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type logger interface {
	Info(text ...any)
	Warn(text ...any)
	Error(text ...any)
}

func WriteSuccessResponse(w http.ResponseWriter, data any) {
	successResponse := SuccessResponse{
		Status: "success",
		Data:   data,
	}

	WriteJSONResponse(w, http.StatusOK, successResponse)
}

func WriteErrorResponse(w http.ResponseWriter, logger logger, err error) {
	var vErr validator.ValidationErrors
	var bErr app_err.BusinessError

	if errors.As(err, &vErr) {
		var errorMessages []string
		for _, fe := range vErr {
			errorMessages = append(errorMessages, cust_validator.GetErrorMsg(fe).Error())
		}
		combinedErrorMessage := strings.Join(errorMessages, ", ")
		errorResponse := ErrorResponse{
			Status: "error",
			Error: Error{
				Code:    bErr.Code(),
				Message: combinedErrorMessage,
			},
		}
		WriteJSONResponse(w, http.StatusBadRequest, errorResponse)
	} else if errors.As(err, &bErr) {
		errorResponse := ErrorResponse{
			Status: "error",
			Error: Error{
				Code:    bErr.Code(),
				Message: bErr.Error(),
			},
		}

		logger.Warn(err)

		WriteJSONResponse(w, http.StatusBadRequest, errorResponse)

	} else {
		errorResponse := ErrorResponse{
			Status: "error",
			Error: Error{
				Code:    "InternalServerError",
				Message: "Что-то пошло не так, попробуйте еще раз",
			},
		}

		logger.Error(err)

		WriteJSONResponse(w, http.StatusInternalServerError, errorResponse)
	}
}

func WriteJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
