package tenant

import (
	"context"
	"strings"

	"github.com/callmhejerry/sms/internal/identity"
	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/auth"
	"github.com/callmhejerry/sms/internal/shared/database"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/callmhejerry/sms/internal/shared/validation"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	queries *store.Queries
	pool    *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool, queries *store.Queries) *Service {
	return &Service{
		pool:    pool,
		queries: queries,
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

	var newTenant store.Tenant
	err := database.WithTx(ctx, service.pool, func(q *store.Queries) error {
		//1.  CREATE TENANT
		t, err := service.queries.CreateTenant(ctx, store.CreateTenantParams{
			Name: name,
			Slug: slug,
		})

		if err != nil {
			return database.TranslateError(err)
			// return nil, apierror.Internal(err, "Failed to create account")
		}
		newTenant = t

		//2. CREATE USER
		hash, err := auth.HashPassword(input.Owner.Password)
		if err != nil {
			return apierror.Internal(err, "failed to hash password")
		}

		newUser, err := service.queries.CreateUser(ctx, store.CreateUserParams{
			TenantID:     newTenant.ID,
			Email:        email,
			FirstName:    firstName,
			LastName:     lastName,
			PasswordHash: hash,
		})
		if err != nil {
			return database.TranslateError(err)
		}

		//3. CREATE OWNER ROLE
		ownerDescription := "Full access to everything within the school"
		ownerRole, err := service.queries.CreateRole(ctx, store.CreateRoleParams{
			TenantID:    newTenant.ID,
			Name:        string(identity.Owner),
			Description: &ownerDescription,
		})

		if err != nil {
			return database.TranslateError(err)
		}

		//4. Assign owner role
		err = service.queries.AssignRoleToUser(ctx, store.AssignRoleToUserParams{
			UserID: newUser.ID,
			RoleID: ownerRole.ID,
		})
		if err != nil {
			return database.TranslateError(err)
		}

		defaultRoles := []struct {
			Name        string
			Description string
		}{
			{Name: string(identity.Admin), Description: "Administrative access"},
			{Name: string(identity.Teacher), Description: "Can manage classes, subject and grades"},
			{Name: string(identity.Accountant), Description: "Can manage fees and payments"},
		}

		for _, r := range defaultRoles {
			_, err := service.queries.CreateRole(ctx, store.CreateRoleParams{
				TenantID:    newTenant.ID,
				Name:        r.Name,
				Description: &r.Description,
			})
			if err != nil {
				return database.TranslateError(err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
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
