package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/KevenMarioN/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func TestEntry(t *testing.T) {
	t.Run("Should be able create a entry", func(t *testing.T) {
		createRandomEntry(t, 0)
	})

	t.Run("Should be able get a new entry", func(t *testing.T) {
		entryNew := createRandomEntry(t, 0)

		entryGet, err := testQueries.GetEntry(context.Background(), entryNew.ID)

		require.NoError(t, err)
		require.NotEmpty(t, entryGet)

		require.Equal(t, entryNew.ID, entryGet.ID)
		require.Equal(t, entryNew.AccountID, entryGet.AccountID)
		require.Equal(t, entryNew.Amount, entryGet.Amount)
		require.WithinDuration(t, entryNew.CreatedAt, entryGet.CreatedAt, time.Second)
	})

	t.Run("Should be not able get entry", func(t *testing.T) {
		entryGet, err := testQueries.GetEntry(context.Background(), util.RandomInt(50, 100))
		require.Error(t, err)
		require.EqualError(t, err, sql.ErrNoRows.Error())
		require.Empty(t, entryGet)
	})

	t.Run("Should be able list of entries", func(t *testing.T) {
		account := createRandomAccount(t)
		for i := 0; i < 10; i++ {
			createRandomEntry(t, account.ID)
		}

		arg := ListEntriesParams{
			Limit:     5,
			Offset:    5,
			AccountID: account.ID,
		}

		entries, err := testQueries.ListEntries(context.Background(), arg)

		require.NoError(t, err)
		require.Len(t, entries, 5)

		for _, entry := range entries {
			require.NotEmpty(t, entry)
		}
	})
}

func createRandomEntry(t *testing.T, id int64) Entry {
	var accountID int64
	if id > 0 {
		accountID = id
	} else {
		account := createRandomAccount(t)
		accountID = account.ID
	}
	arg := CreateEntryParams{
		AccountID: accountID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, accountID)
	require.Equal(t, entry.Amount, arg.Amount)
	require.NotEmpty(t, entry.CreatedAt)
	return entry
}
