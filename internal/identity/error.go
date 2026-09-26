package identity

import (
	"errors"
	"net/http"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUserWithEmailAlreadyExist = apierror.New("user_with_email_already_exists", "A account with this email already exist", http.StatusBadRequest, nil, nil)
	ErrUserNotFound              = apierror.NotFound("User not found")

	ErrInvalidCredentials = apierror.New("invalid_credentials", "Invalid email or password", http.StatusUnauthorized, nil, nil)
	ErrInactiveUser       = apierror.New("user_inactive", "Invalid email or password", http.StatusUnauthorized, nil, nil)
)

func translateUserError(err error) *apierror.AppError {
	if err == nil {
		return nil
	}

	// Expected database/domain errors.
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrUserNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case apierror.UniqueViolation:
			switch pgErr.ConstraintName {
			case "users_email_key":
				return ErrUserWithEmailAlreadyExist
			}
		}
	}

	// Unexpected infrastructure error.
	return apierror.Internal(err, "Something went wrong")
}
