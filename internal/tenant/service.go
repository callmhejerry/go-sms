package tenant

import (
	"context"
	"log/slog"
	"strings"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/auth"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/callmhejerry/sms/internal/shared/validation"
	"github.com/google/uuid"
)

type Service struct {
	queries *store.Queries
	logger  *slog.Logger
}

func NewService(queries *store.Queries, logger *slog.Logger) *Service {
	return &Service{
		queries: queries,
		logger:  logger,
	}
}

func (service *Service) CreateTenant(ctx context.Context, input CreateTenantRequest) (*store.Tenant, error) {
	name := strings.TrimSpace(input.Name)
	slug := strings.TrimSpace(strings.ToLower(input.Slug))
	email := strings.TrimSpace(strings.ToLower(input.Owner.Email))
	firstName := strings.TrimSpace(input.Owner.FirstName)
	lastName := strings.TrimSpace(input.Owner.LastName)

	if name == "" {
		return nil, apierror.Validation("name is required")
	}
	if slug == "" {
		return nil, apierror.Validation("slug is required")
	}
	if email == "" {
		return nil, apierror.Validation("owner email is required")
	}
	if !validation.IsValidEmail(email) {
		return nil, apierror.Validation("invalid owner email")
	}
	if firstName == "" || len(firstName) < 3 {
		return nil, apierror.Validation("first_name is required or invalid")
	}
	if lastName == "" || len(lastName) < 3 {
		return nil, apierror.Validation("last_name is required or invalid")
	}

	if len(input.Owner.Password) < 8 {
		return nil, apierror.Validation("password must be atleast 8 characters")
	}

	//1.  CREATE TENANT
	newTenant, err := service.queries.CreateTenant(ctx, store.CreateTenantParams{
		Name: name,
		Slug: slug,
	})

	if err != nil {
		return nil, apierror.Internal(err, "Failed to create account")
	}

	//2. CREATE USER
	hash, err := auth.HashPassword(input.Owner.Password)
	if err != nil {
		return nil, apierror.Internal(err, "failed to hash password")
	}

	newUser, err := service.queries.CreateUser(ctx, store.CreateUserParams{
		TenantID:     newTenant.ID,
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: hash,
	})
	if err != nil {
		return nil, apierror.Internal(err, "failed to create user")
	}

	//3. CREATE OWNER ROLE
	ownerDescription := "Full access to everything within the school"
	ownerRole, err := service.queries.CreateRole(ctx, store.CreateRoleParams{
		TenantID:    newTenant.ID,
		Name:        "owner",
		Description: &ownerDescription,
	})

	if err != nil {
		return nil, apierror.Internal(err, "failed to create owner role")
	}

	//4. Assign owner role
	err = service.queries.AssignRoleToUser(ctx, store.AssignRoleToUserParams{
		UserID: newUser.ID,
		RoleID: ownerRole.ID,
	})
	if err != nil {
		return nil, apierror.Internal(err, "failed to assign owner role")
	}

	defaultRoles := []struct {
		Name        string
		Description string
	}{
		{Name: "admin", Description: "Administrative access"},
		{Name: "teacher", Description: "Can manage classes, subject and grades"},
		{Name: "accountant", Description: "Can manage fees and payments"},
	}

	for _, r := range defaultRoles {
		_, err := service.queries.CreateRole(ctx, store.CreateRoleParams{
			TenantID:    newTenant.ID,
			Name:        r.Name,
			Description: &r.Description,
		})
		if err != nil {
			if service.logger != nil {
				service.logger.Error("Failed to create default roles")
			}
			continue
		}
	}
	return &newTenant, nil
}

func (service *Service) GetTenantById(ctx context.Context, id uuid.UUID) (*store.Tenant, error) {
	tenant, err := service.GetTenantById(ctx, id)
	if err != nil {
		return nil, apierror.NotFound("Tenant not found")
	}
	return tenant, nil
}

func (service *Service) GetTenantBySlug(ctx context.Context, slug string) (*store.Tenant, error) {
	tenant, err := service.GetTenantBySlug(ctx, slug)

	if err != nil {
		return nil, apierror.NotFound("Tenant not found")
	}
	return tenant, nil
}

func (service *Service) ListTenants(ctx context.Context) ([]store.Tenant, error) {
	tenants, err := service.queries.ListTenants(ctx)

	if err != nil {
		return nil, apierror.Internal(err, "Failed to list tenants")
	}
	return tenants, nil
}
