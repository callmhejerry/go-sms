package billing

import (
	"context"
	"errors"
	"time"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
}

type BillingRepositoryImpl struct {
	queries *store.Queries
}

func NewBillingRepositoryImpl(queries *store.Queries) *BillingRepositoryImpl {
	return &BillingRepositoryImpl{
		queries: queries,
	}
}

func (repo *BillingRepositoryImpl) CreateFeeType(
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

func (repo *BillingRepositoryImpl) ListFeeTypes(
	ctx context.Context,
	tenantId uuid.UUID,
) ([]store.FeeType, *apierror.AppError) {

	fee_types, err := repo.queries.ListFeeTypes(ctx, tenantId)

	if err != nil {
		return nil, translateFeeTypeError(err)
	}

	return fee_types, nil
}

func (repo *BillingRepositoryImpl) CreateFeeStructure(
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

func (repo *BillingRepositoryImpl) ListFeeStructures(
	ctx context.Context,
	tenantId uuid.UUID,
) ([]store.ListFeeStructuresRow, *apierror.AppError) {
	fee_structures, err := repo.queries.ListFeeStructures(ctx, tenantId)

	if err != nil {
		return nil, translateFeeStructuresError(err)
	}
	return fee_structures, nil
}

func (repo *BillingRepositoryImpl) ListFeeStructuresByAcademicSession(
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

func (repo *BillingRepositoryImpl) AssignFeesToStudent(
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

func (repo *BillingRepositoryImpl) ListStudentFees(
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
