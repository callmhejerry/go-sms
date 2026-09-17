package identity

import (
	"context"
	"net/http"
	"strings"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/auth"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	queries    *store.Queries
	jwtManager *auth.JWTManager
}

func NewService(queries *store.Queries, jwtManager *auth.JWTManager) *Service {
	return &Service{
		queries:    queries,
		jwtManager: jwtManager,
	}
}

var (
	InvalidCredentials = apierror.New("invalid_credentials", "Invalid email or password", http.StatusUnauthorized)
	InactiveUser       = apierror.New("user_inactive", "Invalid email or password", http.StatusUnauthorized)
)

func (service *Service) CreateUser(ctx context.Context, input CreateUserRequest) (*store.User, error) {
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

func ConvertToUserResponse(user *store.User) UserResponse {
	return UserResponse{
		ID:        user.ID.String(),
		TenantID:  user.TenantID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt.Time.String(),
		IsActive:  user.IsActive,
	}
}

func (service *Service) Login(ctx context.Context, request LoginRequest) (*LoginResponse, error) {
	email := strings.TrimSpace(strings.ToLower(request.Email))
	slug := strings.TrimSpace(strings.ToLower(request.TenantSlug))

	if slug == "" {
		return nil, apierror.Validation("tenant_slug is required")
	}
	if email == "" {
		return nil, apierror.Validation("email is required")
	}

	tenant, err := service.queries.GetTenantBySlug(ctx, slug)

	if err != nil {
		return nil, InvalidCredentials
	}

	user, err := service.queries.GetUserByEmail(ctx, store.GetUserByEmailParams{
		Email:    email,
		TenantID: tenant.ID,
	})

	if err != nil {
		return nil, InvalidCredentials
	}

	valid, err := auth.CheckPassword(request.Password, user.PasswordHash)
	if err != nil || !valid {
		return nil, InvalidCredentials
	}

	if !user.IsActive {
		return nil, InactiveUser
	}

	userId, err := uuid.Parse(user.ID.String())

	if err != nil {
		return nil, apierror.Internal(err, "Failed to parse user_id")
	}
	tenantId, err := uuid.Parse(tenant.ID.String())

	if err != nil {
		return nil, apierror.Internal(err, "Failed to parse tenant_id")
	}

	token, err := service.jwtManager.Generate(userId, tenantId, email)

	if err != nil {
		return nil, apierror.Internal(err, "Failed to generate token")
	}

	return &LoginResponse{
		Token: token,
		User:  ConvertToUserResponse(&user),
	}, nil
}
