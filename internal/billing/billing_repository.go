package billing

import (
	"context"
	"errors"
	"time"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/database"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BillingRepository interface {
	CreateFeeType(
		ctx context.Context,
		tenantId uuid.UUID,
		name, description string,
		isOptional bool,
	) (*store.FeeType, *apierror.AppError)

	CreateFeeStructure(
		ctx context.Context,
		tenantId, feeTypeId, academicSessionId uuid.UUID,
		classId *uuid.UUID,
		amountKobo int,
		dueDate *time.Time,
	) (*store.FeeStructure, *apierror.AppError)

	ListFeeTypes(
		ctx context.Context,
		tenantId uuid.UUID,
	) ([]store.FeeType, *apierror.AppError)

	ListFeeStructures(
		ctx context.Context,
		tenantId uuid.UUID,
	) ([]store.ListFeeStructuresRow, *apierror.AppError)

	ListFeeStructuresByAcademicSession(
		ctx context.Context,
		tenantId uuid.UUID,
		academicSessionId uuid.UUID,
	) ([]store.ListFeeStructuresBySessionRow, *apierror.AppError)

	AssignFeesToStudent(ctx context.Context, tenantId, studentId uuid.UUID, fee_structures []uuid.UUID) ([]store.StudentFee, *apierror.AppError)

	ListStudentFees(
		ctx context.Context,
		tenantId, studentId uuid.UUID,
	) ([]store.ListStudentFeesRow, *apierror.AppError)

	RecordPayment(
		ctx context.Context,
		tenantId, receivedBy uuid.UUID,
		request RecordPaymentRequest,
	) (*store.Payment, *apierror.AppError)

	ListStudentPayments(
		ctx context.Context,
		tenantId, studentId uuid.UUID,
	) ([]store.Payment, *apierror.AppError)

	GetStudentFeeSummary(
		ctx context.Context, tenantId, studentId uuid.UUID,
	) (*store.GetStudentFeeSummaryRow, *apierror.AppError)

	ListStudentOutstandingFees(
		ctx context.Context,
		tenantId, studentId uuid.UUID,
	) ([]store.ListOutstandingFeesByStudentRow, *apierror.AppError)

	ListOutstandingFees(
		ctx context.Context,
		tenantId uuid.UUID,
	) ([]store.ListOutstandingFeesRow, *apierror.AppError)
}

type billingRepositoryImpl struct {
	queries *store.Queries
	pool    *pgxpool.Pool
}

func NewBillingRepositoryImpl(queries *store.Queries) BillingRepository {
	return &billingRepositoryImpl{
		queries: queries,
	}
}

func (repo *billingRepositoryImpl) CreateFeeType(
	ctx context.Context,
	tenantId uuid.UUID,
	name, description string,
	isOptional bool,
) (*store.FeeType, *apierror.AppError) {
	fee_type, err := repo.queries.CreateFeeType(ctx, store.CreateFeeTypeParams{
		TenantID:    tenantId,
		Name:        name,
		Description: &description,
		IsOptional:  isOptional,
	})

	if err != nil {
		return nil, translateFeeTypeError(err)
	}

	return &fee_type, nil
}

func (repo *billingRepositoryImpl) ListFeeTypes(
	ctx context.Context,
	tenantId uuid.UUID,
) ([]store.FeeType, *apierror.AppError) {

	fee_types, err := repo.queries.ListFeeTypes(ctx, tenantId)

	if err != nil {
		return nil, translateFeeTypeError(err)
	}

	return fee_types, nil
}

func (repo *billingRepositoryImpl) CreateFeeStructure(
	ctx context.Context,
	tenantId, feeTypeId, academicSessionId uuid.UUID,
	classId *uuid.UUID,
	amountKobo int,
	dueDate *time.Time,
) (*store.FeeStructure, *apierror.AppError) {
	fee_structure, err := repo.queries.CreateFeeStructure(ctx, store.CreateFeeStructureParams{
		TenantID:          tenantId,
		FeeTypeID:         feeTypeId,
		AcademicSessionID: academicSessionId,
		ClassID:           classId,
		AmountKobo:        int64(amountKobo),
		DueDate:           dueDate,
	})

	if err != nil {
		return nil, translateFeeStructuresError(err)
	}
	return &fee_structure, nil
}

func (repo *billingRepositoryImpl) ListFeeStructures(
	ctx context.Context,
	tenantId uuid.UUID,
) ([]store.ListFeeStructuresRow, *apierror.AppError) {
	fee_structures, err := repo.queries.ListFeeStructures(ctx, tenantId)

	if err != nil {
		return nil, translateFeeStructuresError(err)
	}
	return fee_structures, nil
}

func (repo *billingRepositoryImpl) ListFeeStructuresByAcademicSession(
	ctx context.Context,
	tenantId uuid.UUID,
	academicSessionId uuid.UUID,
) ([]store.ListFeeStructuresBySessionRow, *apierror.AppError) {
	fee_structures, err := repo.queries.ListFeeStructuresBySession(ctx, store.ListFeeStructuresBySessionParams{
		TenantID:          tenantId,
		AcademicSessionID: academicSessionId,
	})

	if err != nil {
		return nil, translateFeeStructuresError(err)
	}

	return fee_structures, nil
}

func (repo *billingRepositoryImpl) AssignFeesToStudent(
	ctx context.Context,
	tenantId, studentId uuid.UUID,
	fee_structures []uuid.UUID,
) ([]store.StudentFee, *apierror.AppError) {
	var created []store.StudentFee

	for _, fee_structure_id := range fee_structures {
		fee_structure, err := repo.queries.GetFeeStructureById(ctx, store.GetFeeStructureByIdParams{
			ID:       fee_structure_id,
			TenantID: tenantId,
		})

		if err != nil {
			return nil, translateFeeStructuresError(err)
		}

		studentFee, err := repo.queries.CreateStudentFee(
			ctx, store.CreateStudentFeeParams{
				TenantID:       tenantId,
				StudentID:      studentId,
				FeeStructureID: fee_structure.ID,
				AmountKobo:     fee_structure.AmountKobo,
				DueDate:        fee_structure.DueDate,
			},
		)

		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return nil, translateStudentFeesError(err)
			}
		}
		created = append(created, studentFee)
	}

	return created, nil
}

func (repo *billingRepositoryImpl) ListStudentFees(
	ctx context.Context,
	tenantId, studentId uuid.UUID,
) ([]store.ListStudentFeesRow, *apierror.AppError) {
	studentFees, err := repo.queries.ListStudentFees(
		ctx, store.ListStudentFeesParams{
			TenantID:  tenantId,
			StudentID: studentId,
		},
	)
	if err != nil {
		return nil, translateStudentFeesError(err)
	}
	return studentFees, nil
}

func (repo *billingRepositoryImpl) RecordPayment(
	ctx context.Context,
	tenantId, receivedBy uuid.UUID,
	request RecordPaymentRequest,
) (*store.Payment, *apierror.AppError) {
	//TODO : YOU MIGHT NEED TO COME BACK TO THIS
	if request.PaymentMethod == "" {
		request.PaymentMethod = "cash"
	}

	var payment store.Payment

	err := database.WithTx(ctx, repo.pool, func(queries *store.Queries) error {
		p, err := queries.CreatePayment(ctx, store.CreatePaymentParams{
			TenantID:      tenantId,
			StudentID:     request.StudentID,
			AmountKobo:    int64(request.AmountInKobo),
			PaymentMethod: request.PaymentMethod,
			Reference:     request.Reference,
			ReceivedBy:    &receivedBy,
			Notes:         request.Notes,
		})

		if err != nil {
			return translatePaymentsError(err)
		}
		payment = p

		for _, allocation := range request.Allocations {
			_, err := queries.CreatePaymentAllocation(ctx, store.CreatePaymentAllocationParams{
				PaymentID:    p.ID,
				StudentFeeID: allocation.StudentFeeID,
				AmountKobo:   int64(allocation.AmountInKobo),
			})
			if err != nil {
				return translatePaymentAllocationError(err)
			}

			_, err = queries.UpdateStudentFeePayment(
				ctx, store.UpdateStudentFeePaymentParams{
					ID:             allocation.StudentFeeID,
					TenantID:       tenantId,
					AmountPaidKobo: int64(allocation.AmountInKobo),
				},
			)
			if err != nil {
				return translateStudentFeesError(err)
			}
		}

		return nil
	})

	if err != nil {
		var apierr apierror.AppError
		if errors.As(err, &apierr) {
			return nil, &apierr
		}
		return nil, apierror.ErrInternal
	}
	return &payment, nil
}

func (repo *billingRepositoryImpl) ListStudentPayments(
	ctx context.Context,
	tenantId, studentId uuid.UUID,
) ([]store.Payment, *apierror.AppError) {
	payments, err := repo.queries.ListPaymentsByStudent(ctx, store.ListPaymentsByStudentParams{
		TenantID:  tenantId,
		StudentID: studentId,
	})

	if err != nil {
		return nil, translatePaymentsError(err)
	}
	return payments, nil
}

func (repo *billingRepositoryImpl) GetStudentFeeSummary(
	ctx context.Context,
	tenantId, studentId uuid.UUID,
) (*store.GetStudentFeeSummaryRow, *apierror.AppError) {
	summary, err := repo.queries.GetStudentFeeSummary(ctx, store.GetStudentFeeSummaryParams{
		TenantID:  tenantId,
		StudentID: studentId,
	})

	if err != nil {
		return nil, translateStudentFeesError(err)
	}
	return &summary, nil
}

func (repo *billingRepositoryImpl) ListOutstandingFees(
	ctx context.Context,
	tenantId uuid.UUID,
) ([]store.ListOutstandingFeesRow, *apierror.AppError) {

	outstandingFees, err := repo.queries.ListOutstandingFees(ctx, tenantId)
	if err != nil {
		return nil, translateStudentFeesError(err)
	}
	return outstandingFees, nil
}

func (repo *billingRepositoryImpl) ListStudentOutstandingFees(
	ctx context.Context,
	tenantId, studentId uuid.UUID,
) ([]store.ListOutstandingFeesByStudentRow, *apierror.AppError) {
	outstandingFees, err := repo.queries.ListOutstandingFeesByStudent(ctx, store.ListOutstandingFeesByStudentParams{
		TenantID:  tenantId,
		StudentID: studentId,
	})

	if err != nil {
		return nil, translateStudentFeesError(err)
	}
	return outstandingFees, nil
}
