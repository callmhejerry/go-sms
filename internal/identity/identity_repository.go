package identity

import (
	"context"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
)

type IdentityRepository interface {
	CreateNewUser(
		ctx context.Context,
		firstName, lastName, email, password string,
	) (*store.User, *apierror.AppError)

	GetUserById(ctx context.Context, userId uuid.UUID) (*store.User, *apierror.AppError)
	GetUserByEmail(ctx context.Context, email string) (*store.User, *apierror.AppError)
}

type identityRepositoryImpl struct {
	queries *store.Queries
}

func NewRepositoryImpl(
	queries *store.Queries,
) IdentityRepository {
	return &identityRepositoryImpl{
		queries: queries,
	}
}

func (repo *identityRepositoryImpl) CreateNewUser(
	ctx context.Context,
	firstName, lastName, email, passwordHash string,
) (*store.User, *apierror.AppError) {
	newUser, err := repo.queries.CreateUser(ctx, store.CreateUserParams{
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: passwordHash,
	})

	if err != nil {
		return nil, translateUserError(err)
	}
	return &newUser, nil
}

func (repo *identityRepositoryImpl) GetUserById(
	ctx context.Context,
	userId uuid.UUID,
) (*store.User, *apierror.AppError) {
	user, err := repo.queries.GetUserByID(ctx, userId)

	if err != nil {
		return nil, translateUserError(err)
	}

	return &user, nil
}

func (repo *identityRepositoryImpl) GetUserByEmail(
	ctx context.Context,
	email string,
) (*store.User, *apierror.AppError) {
	user, err := repo.queries.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, translateUserError(err)
	}

	return &user, nil
}
