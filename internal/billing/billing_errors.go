package billing

import (
	"errors"
	"net/http"

	academicsession "github.com/callmhejerry/sms/internal/academic/academic_session"
	"github.com/callmhejerry/sms/internal/academic/classes"
	"github.com/callmhejerry/sms/internal/academic/student"
	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrFeeTypeAlreadyExist = apierror.New("fee_type_name_exist", "Fee type with the same name already exist", http.StatusConflict, nil, nil)

	ErrFeeTypeNotFound = apierror.NotFound("Fee type not found")

	ErrFeeStructuresNotFound = apierror.NotFound("Fee structure not found")

	ErrFeeStructureAmountKoboLessThanZero = apierror.New("fee_amount_less_than_zero", "Fee amount must be greater than zero", http.StatusBadRequest, nil, nil)

	ErrFeeStructureAlreadyExist = apierror.New("fee_structure_already_exist", "Fee structure for this fee type , academic session and class already exist", http.StatusConflict, nil, nil)

	ErrStudentFeeAlreadyExist = apierror.New("student_fee_already_exist", "Fee already exist for this student", http.StatusConflict, nil, nil)

	ErrStudentFeeAmountLessThanZero = apierror.New("fee_amount_less_than_zero", "Fee amount must be greater than zero", http.StatusBadRequest, nil, nil)

	ErrStudentFeeNotFound = apierror.NotFound("Student fee not found")

	ErrPaymentNotFound = apierror.NotFound("Payment with id not found")

	ErrPaymentAmountLessThanZero = apierror.New("payment_amount_less_than_zero", "Payment amount must be greater than zero", http.StatusBadRequest, nil, nil)

	ErrPaymentAllocationAlreadyExist = apierror.New(
		"payment_allocation_already_exist", "Payment allocation for this payment already exist",
		http.StatusConflict, nil, nil,
	)

	ErrPaymentAllocationNotFound = apierror.NotFound("Payment allocation not found")

	ErrPaymentAllocationLessThanZero = apierror.New("payment_allocation_less_than_zero", "Payment allocation must be greater than zero", http.StatusBadRequest, nil, nil)
)

func translateFeeTypeError(err error) *apierror.AppError {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrFeeTypeNotFound
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case apierror.UniqueViolation:
			return ErrFeeTypeAlreadyExist
		}
	}

	return apierror.Internal(err, "Something went wrong")
}

func translateFeeStructuresError(err error) *apierror.AppError {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrFeeStructuresNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case apierror.UniqueViolation:
			return ErrFeeStructureAlreadyExist
		case apierror.CheckViolation:
			return ErrFeeStructureAmountKoboLessThanZero
		case apierror.ForeignKeyViolation:
			switch pgErr.ColumnName {
			case "academic_session_id":
				return academicsession.ErrAcademicSessionNotFound
			case "class_id":
				return classes.ErrClassNotFound
			case "fee_type_id":
				return ErrFeeTypeNotFound
			}
		}
	}

	return apierror.Internal(err, "Something went wrong")
}

func translateStudentFeesError(err error) *apierror.AppError {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrStudentFeeNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, *&pgErr) {
		switch pgErr.Code {
		case apierror.ForeignKeyViolation:
			switch pgErr.ColumnName {
			case "student_id":
				return student.ErrStudentNotFound
			case "fee_structure_id":
				return ErrFeeStructuresNotFound
			case "fee_type_id":
				return ErrFeeTypeNotFound
			}
		}
	}
	return apierror.Internal(err, "Something went wrong")
}

func translatePaymentsError(err error) *apierror.AppError {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrPaymentNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case apierror.CheckViolation:
			return ErrPaymentAmountLessThanZero
		case apierror.ForeignKeyViolation:
			switch pgErr.ColumnName {
			case "student_id":
				return student.ErrStudentNotFound
			case "received_by":
				return apierror.NotFound("received_by user not found")
			}
		}
	}
	return apierror.ErrInternal
}

func translatePaymentAllocationError(err error) *apierror.AppError {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrPaymentAllocationNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case apierror.CheckViolation:
			return ErrPaymentAllocationLessThanZero
		case apierror.ForeignKeyViolation:
			switch pgErr.ColumnName {
			case "payment_id":
				return ErrPaymentNotFound
			case "student_fee_id":
				return ErrStudentFeeNotFound
			}
		case apierror.UniqueViolation:
			switch pgErr.ConstraintName {
			case "payment_allocation_unique":
				return ErrPaymentAllocationAlreadyExist
			}
		}
	}
	return apierror.Internal(err, "Something went wrong")
}
