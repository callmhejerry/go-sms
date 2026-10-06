package identity

import (
	"context"
	"testing"
	"time"

	"github.com/callmhejerry/sms/internal/shared/auth"
	"github.com/callmhejerry/sms/internal/shared/database"
	"github.com/callmhejerry/sms/internal/shared/store"
	testsuite "github.com/callmhejerry/sms/internal/shared/test_suite"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshToken_Success(t *testing.T) {
	userID := uuid.New()
	refreshTokenID := uuid.New()
	refreshTokenHash := "refresh-token-hash"

	repo := testsuite.FakeIdentityRepository{
		User: &store.User{
			ID:           userID,
			Email:        "Chinedujeremiah2002@gmail.com",
			PasswordHash: "password-hash",
			FirstName:    "First name",
			LastName:     "Last name",
			IsActive:     true,
			CreatedAt: pgtype.Timestamptz{
				Time:  time.Now(),
				Valid: true,
			},
			UpdatedAt: pgtype.Timestamptz{
				Time:  time.Now(),
				Valid: true,
			},
		},
		RefreshToken: &store.RefreshToken{
			ID:     refreshTokenID,
			UserID: userID,
			ExpiresAt: pgtype.Timestamptz{
				Time:  time.Now().Add(time.Hour),
				Valid: true,
			},
			RevokedAt: pgtype.Timestamptz{
				Valid: false,
			},
			RefreshTokenHash: refreshTokenHash,
		},
	}

	pool := database.NewTestPool(t)
	queries := store.New(pool)
	service := Service{
		queries:      queries,
		jwtManager:   auth.NewJWTManager("test-secret", int(time.Hour)),
		identityRepo: &repo,
	}

	refreshToken, err := service.RefreshToken(context.Background(), "some-refresh-token")

	require.Nil(t, err)
	require.NotNil(t, refreshToken)

	assert.NotEmpty(t, refreshToken.AccessToken)
	assert.NotEmpty(t, refreshToken.RefreshToken)

}
