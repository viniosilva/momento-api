package adapters_test

import (
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"momento/internal/auth/adapters"
	"momento/internal/auth/domain"
)

const (
	testToken  = "test-refresh-token"
	testUserID = "user-123"
	testEmail  = "user@example.com"
	testTTL    = 7 * 24 * time.Hour

	storedPayload = `{"user_id":"user-123","email":"user@example.com"}`

	anyToken = `.+`
)

func TestSecureTokenService_Generate(t *testing.T) {
	t.Run("should generate and save token successfully", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		svc := adapters.NewSecureTokenService(db, testTTL)

		mock.Regexp().ExpectSet(anyToken, anyToken, testTTL).SetVal("OK")

		token, err := svc.Generate(t.Context(), testUserID, testEmail)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when redis Set fails", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		svc := adapters.NewSecureTokenService(db, testTTL)

		mock.Regexp().ExpectSet(anyToken, anyToken, testTTL).SetErr(assert.AnError)

		_, err := svc.Generate(t.Context(), testUserID, testEmail)
		assert.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSecureTokenService_Refresh(t *testing.T) {
	t.Run("should rotate token and return new credentials", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		svc := adapters.NewSecureTokenService(db, testTTL)

		mock.ExpectGet(testToken).SetVal(storedPayload)
		mock.ExpectDel(testToken).SetVal(1)
		mock.Regexp().ExpectSet(anyToken, anyToken, testTTL).SetVal("OK")

		gotUserID, gotEmail, newToken, err := svc.Refresh(t.Context(), testToken)
		require.NoError(t, err)
		assert.Equal(t, testUserID, gotUserID)
		assert.Equal(t, testEmail, gotEmail)
		assert.NotEmpty(t, newToken)
		assert.NotEqual(t, testToken, newToken)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return ErrRefreshTokenNotFound when token does not exist", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		svc := adapters.NewSecureTokenService(db, testTTL)

		mock.ExpectGet(testToken).RedisNil()

		_, _, _, err := svc.Refresh(t.Context(), testToken)
		assert.ErrorIs(t, err, domain.ErrRefreshTokenNotFound)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when redis Get fails", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		svc := adapters.NewSecureTokenService(db, testTTL)

		mock.ExpectGet(testToken).SetErr(assert.AnError)

		_, _, _, err := svc.Refresh(t.Context(), testToken)
		assert.Error(t, err)
		assert.NotErrorIs(t, err, domain.ErrRefreshTokenNotFound)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when redis Del fails", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		svc := adapters.NewSecureTokenService(db, testTTL)

		mock.ExpectGet(testToken).SetVal(storedPayload)
		mock.ExpectDel(testToken).SetErr(assert.AnError)

		_, _, _, err := svc.Refresh(t.Context(), testToken)
		assert.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when redis Set fails on new token", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		svc := adapters.NewSecureTokenService(db, testTTL)

		mock.ExpectGet(testToken).SetVal(storedPayload)
		mock.ExpectDel(testToken).SetVal(1)
		mock.Regexp().ExpectSet(anyToken, anyToken, testTTL).SetErr(assert.AnError)

		_, _, _, err := svc.Refresh(t.Context(), testToken)
		assert.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
