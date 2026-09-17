package tenant

import (
	"context"
	"strings"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
)

type Service struct {
	queries *store.Queries
}

func NewService(queries *store.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

type CreateTenantInput struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (service *Service) CreateTenant(ctx context.Context, input CreateTenantInput) (*store.Tenant, error) {
	name := strings.TrimSpace(input.Name)
	slug := strings.TrimSpace(strings.ToLower(input.Slug))

	if name == "" {
		return nil, apierror.Validation("name is required")
	}
	if slug == "" {
		return nil, apierror.Validation("slug is required")
	}

	row, err := service.queries.CreateTenant(ctx, store.CreateTenantParams{
		Name: name,
		Slug: slug,
	})

	if err != nil {
		return nil, apierror.Internal(err, "Failed to create account")
	}
	return &row, nil
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
	tenants, err := service.ListTenants(ctx)

	if err != nil {
		return nil, apierror.Internal(err, "Failed to list tenants")
	}
	return tenants, nil
}
