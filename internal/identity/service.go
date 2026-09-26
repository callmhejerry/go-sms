package identity

import (
	"context"
	"slices"
	"strings"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/auth"
	"github.com/callmhejerry/sms/internal/shared/constants"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
)

type Service struct {
	queries      *store.Queries
	jwtManager   *auth.JWTManager
	identityRepo IdentityRepository
}

func NewService(
	queries *store.Queries,
	jwtManager *auth.JWTManager,
	identityRepo IdentityRepository,
) *Service {
	return &Service{
		queries:      queries,
		jwtManager:   jwtManager,
		identityRepo: identityRepo,
	}
}

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

	newUser, err := service.identityRepo.CreateNewUser(
		ctx, firstName, lastName, email, hash,
	)

	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (service *Service) GetUserByID(ctx context.Context, id uuid.UUID, tenant_id uuid.UUID) (*store.User, error) {
	user, err := service.queries.GetUserByID(ctx, id)

	if err != nil {
		return nil, apierror.NotFound("User not found")
	}
	return &user, nil
}

func ConvertToUserResponse(user *store.User) UserResponse {
	return UserResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt.Time.String(),
		IsActive:  user.IsActive,
		Email:     user.Email,
	}
}

func (service *Service) Login(ctx context.Context, request LoginRequest) (*LoginResponse, error) {
	email := strings.TrimSpace(strings.ToLower(request.Email))

	if email == "" {
		return nil, apierror.Validation("email is required")
	}

	user, err := service.identityRepo.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	valid, err := auth.CheckPassword(request.Password, user.PasswordHash)
	if err != nil || !valid {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrInactiveUser
	}

	userId, parseErr := uuid.Parse(user.ID.String())

	if parseErr != nil {
		return nil, apierror.Internal(err, "Failed to parse user_id")
	}

	token, jwtErr := service.jwtManager.Generate(userId, email)

	if jwtErr != nil {
		return nil, apierror.Internal(err, "Failed to generate token")
	}

	return &LoginResponse{
		Token: token,
		User:  ConvertToUserResponse(user),
	}, nil
}

func (service *Service) CreateRole(ctx context.Context, tenantId uuid.UUID, input CreateRoleRequest) (*store.Role, error) {
	name := strings.TrimSpace(strings.ToLower(input.Name))
	if name == "" {
		return nil, apierror.Validation("name is required")
	}
	role, err := service.queries.CreateRole(ctx, store.CreateRoleParams{
		TenantID:    tenantId,
		Name:        name,
		Description: &input.Description,
	})

	if err != nil {
		return nil, apierror.Internal(err, "failed to create role")
	}
	return &role, nil
}

func (service *Service) AssignRole(ctx context.Context, userId, roleId uuid.UUID) error {
	err := service.queries.AssignRoleToUser(ctx, store.AssignRoleToUserParams{
		UserID: userId,
		RoleID: roleId,
	})
	if err != nil {
		return apierror.Internal(err, "failed to create role")
	}
	return nil
}

func (service *Service) GetUserRoles(ctx context.Context, userId uuid.UUID) ([]store.Role, error) {

	roles, err := service.queries.GetUserRoles(ctx, userId)

	if err != nil {
		return nil, apierror.Internal(err, "failed to get user roles")
	}
	return roles, nil
}

func (service *Service) UserHasRole(ctx context.Context, userId uuid.UUID, roleName constants.RoleName) (bool, error) {

	userRoles, err := service.queries.GetUserRoles(ctx, userId)

	if err != nil {
		return false, apierror.Internal(err, "failed to check role")
	}

	if slices.ContainsFunc(userRoles, func(role store.Role) bool {
		return role.Name == string(constants.Owner) || role.Name == string(roleName)
	}) {
		return true, nil
	}
	return false, nil
}

// func (service *Service) GetUser(ctx context.Context)
