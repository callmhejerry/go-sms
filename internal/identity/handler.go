package identity

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
	logger  *slog.Logger
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (handler *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apierror.WriteError(w, apierror.Validation("invalid request body"), handler.logger)
		return
	}

	user, err := handler.service.CreateUser(r.Context(), input)
	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ConvertToUserResponse(user))
}

func (handler *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	tenantId, err := uuid.Parse(r.PathValue("tenant_id"))
	if err != nil {
		apierror.WriteError(w, apierror.Validation("Invalid tenant id"), handler.logger)
		return
	}

	userId, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		apierror.WriteError(w, apierror.Validation("Invalid user id"), handler.logger)
		return
	}

	user, err := handler.service.GetUserByID(r.Context(), userId, tenantId)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ConvertToUserResponse(user))
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apierror.WriteError(w, apierror.Validation("Invalid request body"), handler.logger)
		return
	}

	loginResult, err := handler.service.Login(r.Context(), request)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResult)
}
