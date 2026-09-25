package student

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/middleware"
	"github.com/callmhejerry/sms/internal/shared/utils"
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
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, handler.logger)
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

	student, err := handler.service.CreateStudent(r.Context(), tenantId, request)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

func (handler *Handler) GetStudent(w http.ResponseWriter, r *http.Request) {
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	studentId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		apierror.WriteError(w, apierror.Validation("Invalid student_id"), handler.logger)
		return
	}
	student, err := handler.service.GetStudent(r.Context(), tenantId, studentId)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func (handler *Handler) ListStudentsPage(w http.ResponseWriter, r *http.Request) {
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	offsetRequest := utils.ParseOffsetPagination(r)

	students, err := handler.service.ListStudentsPage(
		r.Context(),
		tenantId,
		offsetRequest.Page,
		offsetRequest.PageSize,
	)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func (handler *Handler) GetStudentParents(w http.ResponseWriter, r *http.Request) {
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	studentID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		apierror.WriteError(w, apierror.Validation("invalid student id"), handler.logger)
		return
	}
	parents, err := handler.service.GetStudentParents(r.Context(), tenantId, studentID)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parents)
}

func (handler *Handler) GetStudentProfile(w http.ResponseWriter, r *http.Request) {
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	studentID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		apierror.WriteError(w, apierror.Validation("invalid student id"), handler.logger)
		return
	}

	studentProfile, err := handler.service.GetStudentProfile(r.Context(), tenantId, studentID)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(studentProfile)
}

func (handler *Handler) SearchStudentsPage(w http.ResponseWriter, r *http.Request) {
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	request := parseSearchStudentQuery(r)

	offsetPagination := utils.ParseOffsetPagination(r)

	students, err := handler.service.SearchStudentPage(
		r.Context(),
		tenantId,
		request,
		offsetPagination.Page,
		offsetPagination.PageSize,
	)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func (handler *Handler) SearchStudentsCursor(w http.ResponseWriter, r *http.Request) {
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	request := parseSearchStudentQuery(r)

	cursor := utils.ParseCursorPagination(r)

	students, err := handler.service.SearchStudentCursor(
		r.Context(),
		tenantId,
		request,
		cursor.PageSize,
		DecodeListStudentCursor(cursor.Next),
	)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func (h *Handler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, h.logger)
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

	student, err := h.service.UpdateStudent(r.Context(), tenantId, studentID, input)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}
