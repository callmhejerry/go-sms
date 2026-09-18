package database

import (
	"errors"
	"strings"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TranslateError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return apierror.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503": // foreign_key_violation
			return apierror.Validation("referenced record does not exist")
		case "23502": // not_null_violation
			return apierror.Validation("missing required field")
		}
	}

	return apierror.Internal(err, "database error")
}

func mapUniqueViolation(pgErr *pgconn.PgError) error {
	constraint := pgErr.ConstraintName

	switch {
	case strings.Contains(constraint, "tenants_slug"):
		return apierror.Conflict("a tenant with this slug already exists")
	case strings.Contains(constraint, "users_tenant_id_email"):
		return apierror.Conflict("a user with this email already exists in this school")
	case strings.Contains(constraint, "roles"):
		return apierror.Conflict("a role with this name already exists")
	default:
		return apierror.Conflict("resource already exists")
	}
}
