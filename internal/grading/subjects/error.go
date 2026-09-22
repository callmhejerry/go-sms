package subjects

import (
	"errors"
	"net/http"

	"github.com/callmhejerry/sms/internal/academic/classes"
	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrSubjectCodeAlreadyExist = apierror.New("subject_code_already_exists", "Subject code already exist", http.StatusBadRequest, nil, nil)

	ErrSubjectNotFound = apierror.NotFound("Subject not found")

	ErrSubjectNameAlreadyExist = apierror.New("subject_name_already_exists", "Subject name already exist", http.StatusConflict, nil, nil)

	ErrClassSubjectAlreadyExist = apierror.New("class_subject_already_exists", "This subject already exist for this class", http.StatusConflict, nil, nil)

	ErrTeacherAssignmentAlreadyExist = apierror.New("teacher_assignment_already_exists", "Teacher already exist for this class, subject and academic session", http.StatusConflict, nil, nil)
)

func translateSubjectError(err error) *apierror.AppError {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrSubjectNotFound
	}

	var pgErr pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case apierror.UniqueViolation:
			switch pgErr.ColumnName {
			case "subject_code_unique":
				return ErrSubjectCodeAlreadyExist
			case "subject_name_unique":
				return ErrSubjectNameAlreadyExist
			case "class_subject_unique":
				return ErrClassSubjectAlreadyExist
			case "teacher_assignment_unique":
				return ErrTeacherAssignmentAlreadyExist
			}
		case apierror.ForeignKeyViolation:
			switch pgErr.ColumnName {
			case "class_id":
				return classes.ErrClassNotFound
			case "subject_id":
				return ErrSubjectNotFound
			case "class_arm_id":
				return classes.ErrClassArmNotFound
			}
		}
	}
	return apierror.Internal(err, "something went wrong")
}
