package student

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/middleware"
	"github.com/callmhejerry/sms/internal/shared/validation"
	"github.com/google/uuid"
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

func (handler *Handler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, handler.logger)
		return
	}

	var request CreateStudentRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apierror.WriteError(w, apierror.Validation("Invalid request body"), handler.logger)
		return
	}

	if err := handler.appValidator.ValidateStruct(request); err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	student, err := handler.service.CreateStudent(r.Context(), claims.TenantID, request)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

func (handler *Handler) GetStudent(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, handler.logger)
		return
	}

	studentId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		apierror.WriteError(w, apierror.Validation("Invalid student_id"), handler.logger)
		return
	}
	student, err := handler.service.GetStudent(r.Context(), claims.TenantID, studentId)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func (handler *Handler) ListStudents(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, handler.logger)
		return
	}

	students, err := handler.service.ListStudents(r.Context(), claims.TenantID)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func (handler *Handler) GetStudentParents(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, handler.logger)
		return
	}

	studentID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		apierror.WriteError(w, apierror.Validation("invalid student id"), handler.logger)
		return
	}
	parents, err := handler.service.GetStudentParents(r.Context(), claims.TenantID, studentID)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parents)
}

func (handler *Handler) GetStudentProfile(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, handler.logger)
		return
	}

	studentID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		apierror.WriteError(w, apierror.Validation("invalid student id"), handler.logger)
		return
	}

	studentProfile, err := handler.service.GetStudentProfile(r.Context(), claims.TenantID, studentID)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(studentProfile)
}

func (handler *Handler) SearchStudents(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, handler.logger)
		return
	}

	var request SearchStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apierror.WriteError(w, apierror.Validation("invalid request body"), handler.logger)
		return
	}

	students, err := handler.service.SearchStudent(r.Context(), claims.TenantID, request)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func (h *Handler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}

	studentID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		apierror.WriteError(w, apierror.Validation("invalid student id"), h.logger)
		return
	}

	var input UpdateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apierror.WriteError(w, apierror.Validation("invalid request body"), h.logger)
		return
	}

	student, err := h.service.UpdateStudent(r.Context(), claims.TenantID, studentID, input)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}
