package tenant

import (
	"context"
	"fmt"
	"strings"

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
		return nil, fmt.Errorf("Name is required")
	}
	if slug == "" {
		return nil, fmt.Errorf("Slug is required")
	}

	row, err := service.queries.CreateTenant(ctx, store.CreateTenantParams{
		Name: name,
		Slug: slug,
	})

	if err != nil {
		return nil, fmt.Errorf("Create Tenant: %w", err)
	}
	return &row, nil
}

func (service *Service) GetTenantById(ctx context.Context, id uuid.UUID) (*store.Tenant, error) {
	tenant, err := service.GetTenantById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Get Tenant by id: %w", err)
	}
	return tenant, nil
}

func (service *Service) GetTenantBySlug(ctx context.Context, slug string) (*store.Tenant, error) {
	tenant, err := service.GetTenantBySlug(ctx, slug)

	if err != nil {
		return nil, fmt.Errorf("Get Tenant by slug: %w", err)
	}
	return tenant, nil
}

func (service *Service) ListTenants(ctx context.Context) ([]store.Tenant, error) {
	tenants, err := service.ListTenants(ctx)

	if err != nil {
		return nil, fmt.Errorf("List Tenants: %w", err)
	}
	return tenants, nil
}
