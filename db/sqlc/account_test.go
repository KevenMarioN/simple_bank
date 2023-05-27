package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/KevenMarioN/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.Balance, arg.Balance)
	require.Equal(t, account.Currency, arg.Currency)
	require.Equal(t, account.Owner, arg.Owner)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestAccount(t *testing.T) {
	t.Run("Should be able create a account", func(t *testing.T) {
		createRandomAccount(t)
	})

	t.Run("Should be able get a new account", func(t *testing.T) {
		newAccount := createRandomAccount(t)
		getAccount, err := testQueries.GetAccount(context.Background(), newAccount.ID)

		require.NoError(t, err)
		require.NotEmpty(t, getAccount)

		require.Equal(t, newAccount.ID, getAccount.ID)
		require.Equal(t, newAccount.Owner, getAccount.Owner)
		require.Equal(t, newAccount.Balance, getAccount.Balance)
		require.Equal(t, newAccount.Currency, getAccount.Currency)
		require.WithinDuration(t, newAccount.CreatedAt, getAccount.CreatedAt, time.Second)
	})

	t.Run("Should be able update account", func(t *testing.T) {
		newAccount := createRandomAccount(t)

		arg := UpdateAccountParams{
			ID:      newAccount.ID,
			Balance: util.RandomMoney(),
		}

		updateAccount, err := testQueries.UpdateAccount(context.Background(), arg)

		require.NoError(t, err)
		require.NotEmpty(t, updateAccount)

		require.Equal(t, newAccount.ID, updateAccount.ID)
		require.Equal(t, newAccount.Owner, updateAccount.Owner)
		require.Equal(t, arg.Balance, updateAccount.Balance)
		require.Equal(t, newAccount.Currency, updateAccount.Currency)
		require.WithinDuration(t, newAccount.CreatedAt, updateAccount.CreatedAt, time.Second)
	})

	t.Run("Should be able delete account", func(t *testing.T) {
		account := createRandomAccount(t)

		err := testQueries.DeleteAccount(context.Background(), account.ID)
		require.NoError(t, err)

		accountGet, err := testQueries.GetAccount(context.Background(), account.ID)
		require.Error(t, err)
		require.EqualError(t, err, sql.ErrNoRows.Error())
		require.Empty(t, accountGet)
	})

	t.Run("Should be able list of accounts", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			createRandomAccount(t)
		}

		arg := ListAccountsParams{
			Limit:  5,
			Offset: 5,
		}

		accounts, err := testQueries.ListAccounts(context.Background(), arg)

		require.NoError(t, err)
		require.Len(t, accounts, 5)

		for _, account := range accounts {
			require.NotEmpty(t, account)
		}
	})
}
