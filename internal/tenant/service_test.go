package tenant

import (
	"context"
	"testing"

	"github.com/callmhejerry/sms/internal/shared/database"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

func cleanup(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	_, err := pool.Exec(context.Background(), `
		TRUNCATE TABLE 
			user_roles,
			roles,
			users,
			tenants
		RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		t.Fatalf("failed to clean database: %v", err)
	}
}

func TestCreateTenant_Success(t *testing.T) {
	ctx := context.Background()
	pool := database.NewTestPool(t)
	defer pool.Close()

	cleanup(t, pool)

	queries := store.New(pool)
	service := NewService(pool, queries)

	unique := t.Name()

	input := CreateTenantRequest{
		Name: "Test secondary school",
		Slug: "test-" + unique,
		Owner: CreateOwnerRequest{
			Email:     "owner-" + unique + "@test.com",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "Owner",
		},
	}
	tenant, err := service.CreateTenant(ctx, input)

	if err != nil {
		t.Fatalf("Expected no error got: %v", err)
		return
	}

	if tenant.Name != input.Name {
		t.Fatalf("Expected name %s, got %s", input.Name, tenant.Name)
	}
}

func TestCreateTenant_DuplicateSlug(t *testing.T) {
	ctx := context.Background()
	pool := database.NewTestPool(t)
	defer pool.Close()

	cleanup(t, pool)

	queries := store.New(pool)
	service := NewService(pool, queries)

	slug := "test-duplicates-slug"

	input := CreateTenantRequest{
		Name: "First school",
		Slug: slug,
		Owner: CreateOwnerRequest{
			Email:     "first@test.com",
			Password:  "password123",
			FirstName: "First",
			LastName:  "Test",
		},
	}
	_, err := service.CreateTenant(ctx, input)

	if err != nil {
		t.Fatalf("First create failed: %v", err)
	}

	// try creating with the same slug
	_, err = service.CreateTenant(ctx, input)
	if err == nil {
		t.Fatal("expected error for duplicate slug, got nil")
	}

}

func TestCreateTenant_InvalidEmail(t *testing.T) {
	ctx := context.Background()
	pool := database.NewTestPool(t)
	defer pool.Close()

	cleanup(t, pool)

	queries := store.New(pool)
	service := NewService(pool, queries)

	input := CreateTenantRequest{
		Name: "Test School",
		Slug: "test-school",
		Owner: CreateOwnerRequest{
			Email:     "invalid-email",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
		},
	}

	_, err := service.CreateTenant(ctx, input)
	if err == nil {
		t.Fatal("expected validation error for invalid email")
	}
}
