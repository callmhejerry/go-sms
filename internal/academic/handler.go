package academic

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

func (handler *Handler) CreateClass(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, handler.logger)
		return
	}

	var request CreateClassRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apierror.WriteError(w, apierror.Validation("Invalid request body"), handler.logger)
		return
	}

	if err := handler.appValidator.ValidateStruct(request); err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	class, err := handler.service.CreateClass(r.Context(), claims.TenantID, request)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(class)
}

func (handler *Handler) ListClasses(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, handler.logger)
		return
	}

	classes, err := handler.service.ListClasses(r.Context(), claims.TenantID)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(classes)
}

func (handler *Handler) CreateClassArm(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, handler.logger)
		return
	}

	var request CreateClassArmRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apierror.WriteError(w, apierror.Validation("Invalid request body"), handler.logger)
		return
	}

	if err := handler.appValidator.ValidateStruct(request); err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	classArm, err := handler.service.CreateClassArm(r.Context(), claims.TenantID, request)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(classArm)
}

func (handler *Handler) ListClassArms(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, handler.logger)
		return
	}

	classId, err := uuid.Parse(r.PathValue("class_id"))
	if err != nil {
		apierror.WriteError(w, apierror.Validation("Invalid class id"), handler.logger)
		return
	}
	classArms, err := handler.service.ListClassArms(r.Context(), claims.TenantID, classId)

	if err != nil {
		apierror.WriteError(w, err, handler.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(classArms)
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
