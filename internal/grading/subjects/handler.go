package subjects

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/middleware"
	"github.com/callmhejerry/sms/internal/shared/validation"
)

type SubjectHandler struct {
	subjectService *SubjectService
	logger         *slog.Logger
	appValidator   *validation.AppValidator
}

func NewHandler(service *SubjectService, logger *slog.Logger, appValidator *validation.AppValidator) *SubjectHandler {
	return &SubjectHandler{
		subjectService: service,
		logger:         logger,
		appValidator:   appValidator,
	}
}

func (h *SubjectHandler) CreateSubject(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}

	var input CreateSubjectRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apierror.WriteError(w, apierror.Validation("invalid request body"), h.logger)
		return
	}

	if err := h.appValidator.ValidateStruct(input); err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	subject, err := h.subjectService.CreateSubject(r.Context(), claims.TenantID, input)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subject)
}

func (h *SubjectHandler) ListSubjects(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}

	subjects, err := h.subjectService.ListSubjects(r.Context(), claims.TenantID)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subjects)
}

func (h *SubjectHandler) AddSubjectToClass(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}

	var input AddSubjectToClassRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apierror.WriteError(w, apierror.Validation("invalid request body"), h.logger)
		return
	}

	if err := h.appValidator.ValidateStruct(input); err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	cs, err := h.subjectService.AddSubjectToClass(r.Context(), claims.TenantID, input.ClassID, input.SubjectID)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cs)
}

func (h *SubjectHandler) AssignTeacher(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}

	var input AssignTeacherRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apierror.WriteError(w, apierror.Validation("invalid request body"), h.logger)
		return
	}

	if err := h.appValidator.ValidateStruct(input); err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	assignment, err := h.subjectService.AssignTeacherToSubject(r.Context(), claims.TenantID, input)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(assignment)
}
