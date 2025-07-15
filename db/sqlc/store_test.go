package db

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: int64(account1.ID),
				ToAccountID:   int64(account2.ID),
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, int64(account1.ID), transfer.FromAccountID)
		require.Equal(t, int64(account2.ID), transfer.ToAccountID)
		actualAmount, err := strconv.ParseInt(transfer.Amount, 10, 64)
		require.Equal(t, amount, actualAmount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransferById(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, int64(account1.ID), fromEntry.AccountID)
		actualAmount, err = strconv.ParseInt(fromEntry.Amount, 10, 64)
		require.Equal(t, -amount, actualAmount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntryById(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, int64(account2.ID), toEntry.AccountID)
		actualAmount, err = strconv.ParseInt(toEntry.Amount, 10, 64)
		require.Equal(t, amount, actualAmount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntryById(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//TODO: check accounts balance
	}
}
