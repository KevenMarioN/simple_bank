package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	t.Run("Should be able a transference entrie two accounts", func(t *testing.T) {
		store := NewStore(testDB)
		fromAccount := createRandomAccount(t)
		toAccount := createRandomAccount(t)
		fmt.Println(">> BEFORE:", fromAccount.Balance, toAccount.Balance)

		n := 5
		amount := int64(10)

		errs := make(chan error)
		results := make(chan TransferTxResult)

		for ii := 0; ii < n; ii++ {
			go func() {
				ctx := context.Background()
				result, err := store.TrasferTx(ctx, TransferTxParams{
					FromAccountID: fromAccount.ID,
					ToAccountID:   toAccount.ID,
					Amount:        amount,
				})
				errs <- err
				results <- result
			}()
		}
		existed := make(map[int]bool)
		for ii := 0; ii < n; ii++ {
			err := <-errs
			require.NoError(t, err)

			result := <-results
			require.NotEmpty(t, result)

			// check transfer
			transfer := result.Transfer
			require.NotEmpty(t, transfer)
			require.Equal(t, fromAccount.ID, transfer.FromAccountID)
			require.Equal(t, toAccount.ID, transfer.ToAccountID)
			require.Equal(t, amount, transfer.Amount)
			require.NotZero(t, transfer.ID)
			require.NotZero(t, transfer.CreatedAt)

			_, err = store.GetTransfer(context.Background(), transfer.ID)
			require.NoError(t, err)

			// check entries
			fromEntry := result.FromEntry
			require.NotEmpty(t, fromEntry)
			require.Equal(t, fromAccount.ID, fromEntry.AccountID)
			require.Equal(t, -amount, fromEntry.Amount)
			require.NotZero(t, fromEntry.ID)
			require.NotZero(t, fromEntry.CreatedAt)

			_, err = store.GetEntry(context.Background(), fromEntry.ID)
			require.NoError(t, err)

			// check entries
			toEntry := result.ToEntry
			require.NotEmpty(t, toEntry)
			require.Equal(t, toAccount.ID, toEntry.AccountID)
			require.Equal(t, amount, toEntry.Amount)
			require.NotZero(t, toEntry.ID)
			require.NotZero(t, toEntry.CreatedAt)

			_, err = store.GetEntry(context.Background(), toEntry.ID)
			require.NoError(t, err)

			//check accounts' balance
			fromAccountT := result.FromAccount
			require.NotEmpty(t, fromAccountT)
			require.Equal(t, fromAccount.ID, fromAccountT.ID)

			toAccountT := result.ToAccount
			require.NotEmpty(t, toAccountT)
			require.Equal(t, toAccount.ID, toAccountT.ID)

			// TODO: check accounts´ balance
			fmt.Println(">> TX:", fromAccountT.Balance, toAccountT.Balance)
			diff1 := fromAccount.Balance - fromAccountT.Balance
			diff2 := toAccountT.Balance - toAccount.Balance
			require.Equal(t, diff1, diff2)
			require.True(t, diff1 > 0)
			require.True(t, diff1%amount == 0)

			k := int(diff1 / amount)
			require.True(t, k >= 1 && k <= n)
			require.NotContains(t, existed, k)
			existed[k] = true
		}

		updatedFrom, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
		require.NoError(t, err)

		updatedTo, err := testQueries.GetAccount(context.Background(), toAccount.ID)
		require.NoError(t, err)

		fmt.Println(">> AFTER:", fromAccount.Balance, toAccount.Balance)
		require.Equal(t, fromAccount.Balance-int64(n)*amount, updatedFrom.Balance)
		require.Equal(t, toAccount.Balance+int64(n)*amount, updatedTo.Balance)
	})

	t.Run("Should be able a transference entrie two accounts deadlock", func(t *testing.T) {
		store := NewStore(testDB)
		account1 := createRandomAccount(t)
		account2 := createRandomAccount(t)
		fmt.Println(">> BEFORE:", account1.Balance, account2.Balance)

		n := 10
		amount := int64(10)
		errs := make(chan error)

		for ii := 0; ii < n; ii++ {
			fromAccountID := account1.ID
			toAccountID := account2.ID

			if ii%2 == 1 {
				fromAccountID = account2.ID
				toAccountID = account1.ID
			}
			go func() {
				_, err := store.TrasferTx(context.Background(), TransferTxParams{
					FromAccountID: fromAccountID,
					ToAccountID:   toAccountID,
					Amount:        amount,
				})
				errs <- err
			}()
		}

		for ii := 0; ii < n; ii++ {
			err := <-errs
			require.NoError(t, err)

		}

		updatedFrom, err := testQueries.GetAccount(context.Background(), account1.ID)
		require.NoError(t, err)

		updatedTo, err := testQueries.GetAccount(context.Background(), account2.ID)
		require.NoError(t, err)

		fmt.Println(">> AFTER:", account1.Balance, account2.Balance)
		require.Equal(t, account1.Balance, updatedFrom.Balance)
		require.Equal(t, account2.Balance, updatedTo.Balance)
	})
}
