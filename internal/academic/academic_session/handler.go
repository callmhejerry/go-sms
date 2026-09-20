package academicsession

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/middleware"
	"github.com/callmhejerry/sms/internal/shared/validation"
)

type Handler struct {
	logger       *slog.Logger
	service      *Service
	appValidator *validation.AppValidator
}

func NewHandler(service *Service, logger *slog.Logger, appValidator *validation.AppValidator) *Handler {
	return &Handler{
		logger:       logger,
		service:      service,
		appValidator: appValidator,
	}
}

func (handler *Handler) CreateAcademicSession(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, handler.logger)
		return
	}

	var request CreateAcademicSessionRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apierror.WriteError(w, apierror.Validation("Invalid request body"), handler.logger)
		return
	}

	if err := handler.appValidator.ValidateStruct(request); err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	session, err := handler.service.CreateAcademicSession(r.Context(), claims.TenantID, request)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

func (handler *Handler) ListAcademicSessions(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, handler.logger)
		return
	}

	academicSessions, err := handler.service.ListAcademicSession(r.Context(), claims.TenantID)
	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(academicSessions)
}

func (handler *Handler) GetCurrentSession(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, handler.logger)
		return
	}

	currentAcademicSession, err := handler.service.GetCurrentAcademicSession(r.Context(), claims.TenantID)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentAcademicSession)
}
