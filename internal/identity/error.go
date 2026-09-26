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
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrUserNotFound
	}

	var pgErr pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case apierror.UniqueViolation:
			if pgErr.ColumnName == "email" {
				return ErrUserWithEmailAlreadyExist
			}
		}
	}
	return apierror.Internal(err, "Something went wrong")
}
