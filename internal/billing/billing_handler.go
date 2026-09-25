package billing

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/middleware"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/callmhejerry/sms/internal/shared/validation"
	"github.com/google/uuid"
)

type BillingHandler struct {
	billingService *BillingService
	logger         *slog.Logger
	validator      *validation.AppValidator
}

func NewBillingHandler(billingService *BillingService, logger *slog.Logger, validator *validation.AppValidator) *BillingHandler {
	return &BillingHandler{
		billingService: billingService,
		logger:         logger,
		validator:      validator,
	}
}

func (h *BillingHandler) CreateFeeType(w http.ResponseWriter, r *http.Request) {
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, h.logger)
		return
	}

	var request CreateFeeTypeRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apierror.WriteError(w, apierror.Validation("Invalid request body"), h.logger)
		return
	}

	if err := h.validator.ValidateStruct(request); err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	feeType, err := h.billingService.CreateFeeType(
		r.Context(),
		tenantId,
		request,
	)

	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(feeType)
}

func (h *BillingHandler) CreateFeeStructure(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, h.logger)
		return
	}

	var request CreatFeeStructureRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apierror.WriteError(w, apierror.Validation("Invalid request body"), h.logger)
		return
	}

	if err := h.validator.ValidateStruct(request); err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	feeStructure, err := h.billingService.CreateFeeStructure(r.Context(), tenantId, request)

	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(feeStructure)
}

func (h *BillingHandler) ListFeeTypes(w http.ResponseWriter, r *http.Request) {
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, h.logger)
		return
	}

	feeTypes, err := h.billingService.ListFeeTypes(r.Context(), tenantId)

	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feeTypes)
}

func (h *BillingHandler) ListFeeStructures(w http.ResponseWriter, r *http.Request) {
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, h.logger)
		return
	}

	var academicSessionId *uuid.UUID
	query := r.URL.Query()

	if value := query.Get("academic_session_id"); value != "" {
		if parsedUUID, err := uuid.Parse(value); err == nil {
			academicSessionId = &parsedUUID
		}
	}

	var feeStructures []store.FeeStructure

	if academicSessionId == nil {
		feeStructures, err = h.billingService.ListFeeStructures(r.Context(), tenantId)
	} else {
		feeStructures, err = h.billingService.ListFeeStructuresByAdmissionId(r.Context(), tenantId, *academicSessionId)
	}

	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feeStructures)
}

func (h *BillingHandler) AssignFeesToStudent(w http.ResponseWriter, r *http.Request) {
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, h.logger)
		return
	}

	var request AssignFeesToStudentRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apierror.WriteError(w, apierror.Validation("Invalid request body"), h.logger)
		return
	}

	if err := h.validator.ValidateStruct(request); err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	studentFees, err := h.billingService.AssignFeesToStudent(r.Context(), tenantId, request)

	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(studentFees)
}

func (h *BillingHandler) ListStudentFees(w http.ResponseWriter, r *http.Request) {
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, h.logger)
		return
	}

	studentID, err := uuid.Parse(r.PathValue("student_id"))
	if err != nil {
		apierror.WriteError(w, apierror.Validation("invalid student id"), h.logger)
		return
	}

	fees, err := h.billingService.ListStudentFees(r.Context(), tenantId, studentID)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fees)
}

func (h *BillingHandler) RecordPayment(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apierror.WriteError(w, apierror.ErrUnauthorized, h.logger)
		return
	}

	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, h.logger)
		return
	}

	var request RecordPaymentRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apierror.WriteError(w, apierror.Validation("Invalid request body"), h.logger)
		return
	}
	if err := h.validator.ValidateStruct(request); err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	payment, err := h.billingService.RecordPayments(
		r.Context(), tenantId, claims.UserID,
		request,
	)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment)
}

func (h *BillingHandler) ListStudentPayments(w http.ResponseWriter, r *http.Request) {
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, h.logger)
		return
	}

	studentID, err := uuid.Parse(r.PathValue("student_id"))
	if err != nil {
		apierror.WriteError(w, apierror.Validation("invalid student id"), h.logger)
		return
	}

	payments, err := h.billingService.ListStudentPayments(r.Context(), tenantId, studentID)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

func (h *BillingHandler) GetStudentFeeSummary(w http.ResponseWriter, r *http.Request) {
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, h.logger)
		return
	}

	studentID, err := uuid.Parse(r.PathValue("student_id"))
	if err != nil {
		apierror.WriteError(w, apierror.Validation("invalid student id"), h.logger)
		return
	}

	summary, err := h.billingService.GetStudentFeeSummary(r.Context(), tenantId, studentID)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

func (h *BillingHandler) ListOutstandingFees(w http.ResponseWriter, r *http.Request) {
	tenantId, found, err := middleware.GetTenantId(r.Context())
	if !found {
		apierror.WriteError(w, err, h.logger)
		return
	}

	fees, err := h.billingService.ListOutstandingFees(r.Context(), tenantId)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fees)
}
