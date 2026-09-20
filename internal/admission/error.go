package admission

import (
	"errors"
	"strings"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/database"
	"github.com/jackc/pgx/v5/pgconn"
)

func translateAdmissionError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" { // unique_violation
			switch {
			case strings.Contains(pgErr.ConstraintName, "academic_sessions"):
				return apierror.Conflict("an academic session with this name already exists")
			case strings.Contains(pgErr.ConstraintName, "classes"):
				return apierror.Conflict("a class with this name already exists")
			case strings.Contains(pgErr.ConstraintName, "class_arms"):
				return apierror.Conflict("an arm with this name already exists for the class")
			}
		}
	}

	// Fall back to generic translation
	return database.TranslateError(err)
}
