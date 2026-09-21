package academicsession

import (
	"errors"
	"net/http"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrAcademicSessionAlreadyExists = apierror.New("academic_session_already_exist", "Academic session with the same name already exist", http.StatusBadRequest, nil, nil)

	ErrAcademicSessionNotFound = apierror.NotFound("Academic session not found")
)

func translateError(err error) *apierror.AppError {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrAcademicSessionNotFound
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return mapUniqueKeyConstraint(pgErr.ConstraintName)
		}
	}

	return apierror.Internal(err, "Something went wrong")
}

func mapUniqueKeyConstraint(constraintName string) *apierror.AppError {
	switch constraintName {
	case "academic_sessions_tenant_id_name_unique":
		return ErrAcademicSessionAlreadyExists
	default:
		return apierror.ErrConflict
	}
}
