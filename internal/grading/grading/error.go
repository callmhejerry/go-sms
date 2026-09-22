package grading

import (
	"errors"
	"net/http"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrMaxScoreMustBeGreaterThanZero = apierror.New("maxscore_must_be_greater_than_zero", "Max scrore must be greater than zero", http.StatusBadRequest, nil, nil)

	ErrWeightMustBeGreaterThanZero = apierror.New("weight_must_be_greater_than_zero", "Weight must be greater than zero", http.StatusBadRequest, nil, nil)

	ErrAssessmentTypeAlreadyExist = apierror.New("assessment_type_already_exists", "Assessment type with the same name already exist", http.StatusConflict, nil, nil)

	ErrAssessmentTypeNotFound = apierror.NotFound("assessment type not found")
)

func translateAssessmentTypeError(err error) *apierror.AppError {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrAssessmentTypeNotFound
	}

	var pgErr pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case apierror.CheckViolation:
			switch pgErr.ColumnName {
			case "max_score":
				return ErrMaxScoreMustBeGreaterThanZero
			case "weight":
				return ErrWeightMustBeGreaterThanZero
			}
		case apierror.UniqueViolation:
			if pgErr.ColumnName == "name" {
				return ErrAssessmentTypeAlreadyExist
			}
		}
	}
	return apierror.Internal(err, "something went wrong")
}
