package identity

import (
	"context"
	"errors"
	"time"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type IdentityRepository interface {
	CreateNewUser(
		ctx context.Context,
		firstName, lastName, email, password string,
	) (*store.User, *apierror.AppError)

	GetUserById(ctx context.Context, userId uuid.UUID) (*store.User, *apierror.AppError)
	GetUserByEmail(ctx context.Context, email string) (*store.User, *apierror.AppError)
	CreateRefreshToken(ctx context.Context, userId uuid.UUID, refreshTokenHash string, expiresAt time.Duration) *apierror.AppError
	GetRefreshToken(ctx context.Context, refreshTokenHash string) (*store.RefreshToken, *apierror.AppError)
	RevokeRefreshToken(ctx context.Context, id uuid.UUID) *apierror.AppError
	DeleteRefreshToken(ctx context.Context, id uuid.UUID) *apierror.AppError
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

func (repo *identityRepositoryImpl) CreateRefreshToken(
	ctx context.Context,
	userId uuid.UUID,
	refreshTokenHash string,
	expiresAt time.Duration,
) *apierror.AppError {
	_, err := repo.queries.CreateRefreshToken(
		ctx, store.CreateRefreshTokenParams{
			RefreshTokenHash: refreshTokenHash,
			UserID:           userId,
			ExpiresAt: pgtype.Timestamptz{
				Time:  time.Now().Add(expiresAt),
				Valid: true,
			},
		},
	)
	if err != nil {
		return apierror.Internal(err, "Something went wrong")
	}
	return nil
}

func (repo *identityRepositoryImpl) RevokeRefreshToken(
	ctx context.Context,
	id uuid.UUID,
) *apierror.AppError {
	err := repo.queries.RevokeRefreshToken(ctx, id)
	if err != nil {
		return apierror.Internal(err, "Something went wrong")
	}
	return nil
}

func (repo *identityRepositoryImpl) GetRefreshToken(
	ctx context.Context,
	refreshTokenHash string,
) (*store.RefreshToken, *apierror.AppError) {
	refreshToken, err := repo.queries.GetRefreshToken(ctx, refreshTokenHash)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NotFound("Refresh token not found")
		}
		return nil, apierror.Internal(err, "something went wrong")
	}

	return &refreshToken, nil
}

func (repo *identityRepositoryImpl) DeleteRefreshToken(
	ctx context.Context,
	id uuid.UUID,
) *apierror.AppError {
	err := repo.queries.DeleteRefreshToken(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return apierror.Internal(err, "something went wrong")
	}

	return nil
}
