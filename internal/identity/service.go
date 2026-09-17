package identity

import (
	"context"
	"strings"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/auth"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	queries *store.Queries
}

func NewService(queries *store.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

type CreateUserInput struct {
	TenantID  uuid.UUID `json:"tenant_id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

func (service *Service) CreateUser(ctx context.Context, input CreateUserInput) (*store.User, error) {
	email := strings.TrimSpace(input.Email)
	firstName := strings.TrimSpace(input.FirstName)
	lastName := strings.TrimSpace(input.LastName)

	if email == "" {
		return nil, apierror.Validation("email is required")
	}
	if input.Password == "" {
		return nil, apierror.Validation("password is required")
	}
	if len(input.Password) < 8 {
		return nil, apierror.Validation("password must be atleast 8 characters")
	}
	if firstName == "" {
		return nil, apierror.Validation("first_name is required")
	}
	if lastName == "" {
		return nil, apierror.Validation("last_name is required")
	}

	hash, err := auth.HashPassword(input.Password)

	if err != nil {
		return nil, apierror.Internal(err, "failed to hash password")
	}

	newUser, err := service.queries.CreateUser(ctx, store.CreateUserParams{
		TenantID: pgtype.UUID{
			Bytes: [16]byte(input.TenantID),
			Valid: true,
		},
		Email:        input.Email,
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: hash,
	})

	if err != nil {
		return nil, apierror.Internal(err, "Failed to create new user")
	}
	return &newUser, nil
}

func (service *Service) GetUserByID(ctx context.Context, id uuid.UUID, tenant_id uuid.UUID) (*store.User, error) {
	user, err := service.queries.GetUserByID(ctx, store.GetUserByIDParams{
		ID: pgtype.UUID{
			Bytes: id,
			Valid: true,
		},
		TenantID: pgtype.UUID{
			Bytes: tenant_id,
			Valid: true,
		},
	})

	if err != nil {
		return nil, apierror.NotFound("User not found")
	}
	return &user, nil
}
