package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	t.Run("Should be able create transfer", func(t *testing.T) {
		toAccount := createRandomAccount(t)
		fromAccount := createRandomAccount(t)
		createRandomTransfer(t, fromAccount.ID, toAccount.ID)
	})

	t.Run("Should be able get transfer", func(t *testing.T) {
		toAccount := createRandomAccount(t)
		fromAccount := createRandomAccount(t)
		transferNew := createRandomTransfer(t, fromAccount.ID, toAccount.ID)

		transferGet, err := testQueries.GetTransfer(context.Background(), transferNew.ID)

		require.NoError(t, err)
		require.NotEmpty(t, transferGet)

		require.Equal(t, transferNew.ID, transferGet.ID)
		require.Equal(t, transferNew.Amount, transferGet.Amount)
		require.Equal(t, transferNew.FromAccountID, transferGet.FromAccountID)
		require.Equal(t, transferNew.ToAccountID, transferGet.ToAccountID)
		require.WithinDuration(t, transferNew.CreatedAt, transferGet.CreatedAt, time.Second)
	})

	t.Run("Should be able list transfer", func(t *testing.T) {
		toAccount := createRandomAccount(t)
		fromAccount := createRandomAccount(t)

		for ii := 0; ii < 10; ii++ {
			createRandomTransfer(t, fromAccount.ID, toAccount.ID)
			createRandomTransfer(t, toAccount.ID, fromAccount.ID)
		}
		arg := ListTransfersParams{
			FromAccountID: fromAccount.ID,
			ToAccountID:   toAccount.ID,
			Limit:         5,
			Offset:        5,
		}

		transfers, err := testQueries.ListTransfers(context.Background(), arg)
		require.NoError(t, err)
		require.Len(t, transfers, 5)

		for _, transfer := range transfers {
			require.NotEmpty(t, transfer)
			require.True(t, transfer.FromAccountID == fromAccount.ID || transfer.ToAccountID == toAccount.ID)
		}
	})
}

func createRandomTransfer(t *testing.T, accountFromID, toAccountID int64) Transfer {

	arg := CreateTransferParams{
		FromAccountID: accountFromID,
		ToAccountID:   toAccountID,
		Amount:        120,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, accountFromID)
	require.Equal(t, transfer.ToAccountID, toAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotEmpty(t, transfer.CreatedAt)
	require.NotEmpty(t, transfer.ID)

	return transfer
}
