package testsuite

import (
	"context"
	"time"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
)

type FakeIdentityRepository struct {
	RefreshToken *store.RefreshToken
	User         *store.User

	CreatedRefreshToken bool
	DeletedTokenID      uuid.UUID
}

func (f *FakeIdentityRepository) CreateNewUser(
	ctx context.Context,
	firstName, lastName, email, password string,
) (*store.User, *apierror.AppError) {
	return f.User, nil
}

func (f *FakeIdentityRepository) GetUserById(ctx context.Context, userId uuid.UUID) (*store.User, *apierror.AppError) {
	return f.User, nil
}

func (f *FakeIdentityRepository) GetUserByEmail(ctx context.Context, email string) (*store.User, *apierror.AppError) {
	return f.User, nil
}

func (f *FakeIdentityRepository) CreateRefreshToken(ctx context.Context, userId uuid.UUID, refreshTokenHash string, expiresAt time.Duration) *apierror.AppError {
	f.CreatedRefreshToken = true
	return nil
}

func (f *FakeIdentityRepository) GetRefreshToken(ctx context.Context, refreshTokenHash string) (*store.RefreshToken, *apierror.AppError) {
	return f.RefreshToken, nil
}

func (f *FakeIdentityRepository) RevokeRefreshToken(ctx context.Context, id uuid.UUID) *apierror.AppError {
	return nil
}

func (f *FakeIdentityRepository) DeleteRefreshToken(
	ctx context.Context,
	id uuid.UUID,
) *apierror.AppError {
	f.DeletedTokenID = id
	return nil
}
