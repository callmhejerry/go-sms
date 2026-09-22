package grading

import (
	"errors"
	"net/http"

	academicsession "github.com/callmhejerry/sms/internal/academic/academic_session"
	"github.com/callmhejerry/sms/internal/academic/classes"
	"github.com/callmhejerry/sms/internal/academic/student"
	"github.com/callmhejerry/sms/internal/grading/subjects"
	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrMaxScoreMustBeGreaterThanZero = apierror.New("maxscore_must_be_greater_than_zero", "Max scrore must be greater than zero", http.StatusBadRequest, nil, nil)

	ErrWeightMustBeGreaterThanZero = apierror.New("weight_must_be_greater_than_zero", "Weight must be greater than zero", http.StatusBadRequest, nil, nil)

	ErrAssessmentTypeAlreadyExist = apierror.New("assessment_type_already_exists", "Assessment type with the same name already exist", http.StatusConflict, nil, nil)

	ErrAssessmentTypeNotFound = apierror.NotFound("assessment type not found")

	ErrScoreNotFound = apierror.NotFound("Score not found")

	ErrScoreAlreadyExist = apierror.New("score_already_exist", "Score already exist for this student", http.StatusConflict, nil, nil)

	ErrScoreLessThanZero = apierror.New("score_less_than_zero", "Score must be greater than zero", http.StatusBadRequest, nil, nil)

	ErrScoreGreaterThanMaxScore = apierror.New("score_greater_than_max_score", "Score cannot be greater than the maximum score for this assessment type", http.StatusBadRequest, nil, nil)
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

func translateScoresTypeError(err error) *apierror.AppError {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrScoreNotFound
	}

	var pgErr pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case apierror.CheckViolation:
			switch pgErr.ColumnName {
			case "score":
				return ErrScoreLessThanZero
			}
		case apierror.UniqueViolation:
			if pgErr.ColumnName == "score_unique" {
				return ErrScoreAlreadyExist
			}
		case apierror.ForeignKeyViolation:
			switch pgErr.ColumnName {
			case "student_id":
				return student.ErrStudentNotFound
			case "subject_id":
				return subjects.ErrSubjectNotFound
			case "assessment_type_id":
				return ErrAssessmentTypeNotFound
			case "class_arm_id":
				return classes.ErrClassArmNotFound
			case "academic_session_id":
				return academicsession.ErrAcademicSessionNotFound
			}
		}

	}
	return apierror.Internal(err, "something went wrong")
}
