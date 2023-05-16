package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	t.Run("Should be able a transference entrie two accounts", func(t *testing.T) {
		store := NewStore(testDB)
		fromAccount := createRandomAccount(t)
		toAccount := createRandomAccount(t)

		n := 5
		amount := int64(10)

		errs := make(chan error)
		results := make(chan TransferTxResult)

		for ii := 0; ii < n; ii++ {
			go func() {
				result, err := store.TrasferTx(context.Background(), TransferTxParams{
					FromAccountID: fromAccount.ID,
					ToAccountID:   toAccount.ID,
					Amount:        amount,
				})
				errs <- err
				results <- result
			}()
		}

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

			// TODO: check accounts´ balance
		}
	})
}
