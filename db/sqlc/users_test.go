package db

import (
	"context"
	"testing"
	"time"

	"github.com/KevenMarioN/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: util.RandomString(20),
		Email:          util.RandomString(20),
		FullName:       util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, user.Email, arg.Email)
	require.Equal(t, user.FullName, arg.FullName)
	require.Equal(t, user.HashedPassword, arg.HashedPassword)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestUser(t *testing.T) {
	t.Run("Should be able create a account", func(t *testing.T) {
		createRandomUser(t)
	})

	t.Run("Should be able get a new account", func(t *testing.T) {
		newUser := createRandomUser(t)
		getUser, err := testQueries.GetUser(context.Background(), newUser.Username)

		require.NoError(t, err)
		require.NotEmpty(t, getUser)

		require.Equal(t, newUser.Username, getUser.Username)
		require.Equal(t, newUser.Email, getUser.Email)
		require.Equal(t, newUser.FullName, getUser.FullName)
		require.Equal(t, newUser.PasswordChangedAt, getUser.PasswordChangedAt)
		require.Equal(t, newUser.HashedPassword, getUser.HashedPassword)

		require.WithinDuration(t, newUser.CreatedAt, getUser.CreatedAt, time.Second)
	})

}
