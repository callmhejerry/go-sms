package apierror

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type errorResponse struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func WriteError(w http.ResponseWriter, err error, logger *slog.Logger) {
	var appError *AppError

	if errors.As(err, &appError) {
		writeJSON(w, appError.HTTPStatus, errorResponse{
			Error: errorBody{
				Code:    appError.Code,
				Message: appError.Message,
			},
		})
		return
	}

	if logger != nil {
		logger.Error("unhandled error", slog.String("error", err.Error()))
	}
	writeJSON(w, http.StatusInternalServerError, errorResponse{
		Error: errorBody{
			Code:    ErrInternal.Code,
			Message: ErrInternal.Message,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
