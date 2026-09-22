package billing

import (
	"context"
	"time"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
)

type BillingService struct {
	billingRepository BillingRepository
}

func NewBillingService(billingRepository BillingRepository) *BillingService {
	return &BillingService{
		billingRepository: billingRepository,
	}
}

func (service *BillingService) CreateFeeType(
	ctx context.Context,
	tenantId uuid.UUID,
	request CreateFeeTypeRequest,
) (*store.FeeType, *apierror.AppError) {
	return service.billingRepository.CreateFeeType(ctx, tenantId, request.Name, *request.Description, request.IsOptional)
}

func (service *BillingService) CreateFeeStructure(
	ctx context.Context,
	tenantId uuid.UUID,
	request CreatFeeStructureRequest,
) (*store.FeeStructure, *apierror.AppError) {
	var dueDate *time.Time

	if request.DueDate != nil {
		if parsedDate, err := time.Parse("2026-01-30", *request.DueDate); err == nil {
			dueDate = &parsedDate
		}
	}

	return service.billingRepository.CreateFeeStructure(
		ctx,
		tenantId,
		request.FeeTypeId,
		request.AcademicSessionId, request.ClassId, request.AmountInKobo, dueDate,
	)
}

func (service *BillingService) ListFeeTypes(
	ctx context.Context,
	tenantId uuid.UUID,
) ([]store.FeeType, *apierror.AppError) {
	return service.billingRepository.ListFeeTypes(ctx, tenantId)
}

func (service *BillingService) ListFeeStructures(
	ctx context.Context,
	tenantId uuid.UUID,
) ([]store.FeeStructure, *apierror.AppError) {
	rows, err := service.billingRepository.ListFeeStructures(ctx, tenantId)

	if err != nil {
		return nil, err
	}

	feeStructures := make([]store.FeeStructure, len(rows))

	for _, row := range rows {
		feeStructures = append(feeStructures, store.FeeStructure{
			ID:                row.ID,
			TenantID:          row.TenantID,
			FeeTypeID:         row.FeeTypeID,
			AcademicSessionID: row.AcademicSessionID,
			ClassID:           row.ClassID,
			AmountKobo:        row.AmountKobo,
			DueDate:           row.DueDate,
			CreatedAt:         row.CreatedAt,
			UpdatedAt:         row.UpdatedAt,
		})
	}

	return feeStructures, nil
}

func (service *BillingService) ListFeeStructuresByAdmissionId(
	ctx context.Context,
	tenantId,
	academicSessionId uuid.UUID,
) ([]store.FeeStructure, *apierror.AppError) {
	rows, err := service.billingRepository.ListFeeStructuresByAcademicSession(ctx, tenantId, academicSessionId)

	if err != nil {
		return nil, err
	}

	feeStructures := make([]store.FeeStructure, len(rows))

	for _, row := range rows {
		feeStructures = append(feeStructures, store.FeeStructure{
			ID:                row.ID,
			TenantID:          row.TenantID,
			FeeTypeID:         row.FeeTypeID,
			AcademicSessionID: row.AcademicSessionID,
			ClassID:           row.ClassID,
			AmountKobo:        row.AmountKobo,
			DueDate:           row.DueDate,
			CreatedAt:         row.CreatedAt,
			UpdatedAt:         row.UpdatedAt,
		})
	}

	return feeStructures, nil
}

func (service *BillingService) AssignFeesToStudent(
	ctx context.Context,
	tenantId uuid.UUID,
	request AssignFeesToStudentRequest,
) ([]store.StudentFee, *apierror.AppError) {
	return service.billingRepository.AssignFeesToStudent(ctx, tenantId, request.StudentID, request.FeeStructureIDs)
}

func (service *BillingService) ListStudentFees(
	ctx context.Context,
	tenantId, studentId uuid.UUID,
) ([]store.ListStudentFeesRow, *apierror.AppError) {
	return service.billingRepository.ListStudentFees(ctx, tenantId, studentId)
}

func (service *BillingService) RecordPayments(
	ctx context.Context,
	tenantId, receivedBy uuid.UUID,
	request RecordPaymentRequest,
) (*store.Payment, *apierror.AppError) {
	return service.billingRepository.RecordPayment(
		ctx, tenantId, receivedBy,
		request,
	)
}

func (service *BillingService) ListStudentPayments(
	ctx context.Context, tenantId, studentId uuid.UUID,
) ([]store.Payment, *apierror.AppError) {
	return service.billingRepository.ListStudentPayments(
		ctx, tenantId,
		studentId,
	)
}

func (service *BillingService) GetStudentFeeSummary(
	ctx context.Context,
	tenantId, studentId uuid.UUID,
) (*store.GetStudentFeeSummaryRow, *apierror.AppError) {
	return service.billingRepository.GetStudentFeeSummary(ctx, tenantId, studentId)
}

func (service *BillingService) ListOutstandingFees(
	ctx context.Context, tenantId uuid.UUID,
) ([]store.ListOutstandingFeesRow, *apierror.AppError) {
	return service.billingRepository.ListOutstandingFees(ctx, tenantId)
}

func (service *BillingService) ListStudentOutstandingFees(
	ctx context.Context,
	tenantId, studentId uuid.UUID,
) ([]store.ListOutstandingFeesByStudentRow, *apierror.AppError) {
	return service.billingRepository.ListStudentOutstandingFees(ctx, tenantId, studentId)
}
