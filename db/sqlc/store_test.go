package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)
	n := 5
	amount := decimal.NewFromInt(10)

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

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, int64(account1.ID), transfer.FromAccountID)
		require.Equal(t, int64(account2.ID), transfer.ToAccountID)
		actualAmount := transfer.Amount
		require.Equal(t, amount, actualAmount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransferById(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, int64(account1.ID), fromEntry.AccountID)
		actualAmount = fromEntry.Amount
		require.Equal(t, amount.Neg(), actualAmount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntryById(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, int64(account2.ID), toEntry.AccountID)
		actualAmount = toEntry.Amount
		require.Equal(t, amount, actualAmount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntryById(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)
		fmt.Println(">> before:", account1.Balance, account2.Balance)

		diff1 := account1.Balance.Sub(fromAccount.Balance)
		diff2 := toAccount.Balance.Sub(account2.Balance)

		require.Equal(t, diff1, diff2)
		require.True(t, diff1.Cmp(decimal.NewFromInt(0)) == 1)
		require.True(t, diff1.Mod(amount).IsZero())

		k := diff1.Div(amount)
		require.True(t, k.Cmp(decimal.NewFromInt(1)) >= 0 && k.Cmp(decimal.NewFromInt(int64(n))) <= 0)
		require.NotContains(t, existed, k)
		existed[int(k.IntPart())] = true

	}

	updateAccount1, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccountById(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updateAccount1.Balance, updateAccount2.Balance)

	require.True(t, account1.Balance.Sub(amount.Mul(decimal.NewFromInt(int64(n)))).Equal(updateAccount1.Balance))
	require.True(t, account2.Balance.Add(amount.Mul(decimal.NewFromInt(int64(n)))).Equal(updateAccount2.Balance))

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)
	n := 10
	amount := decimal.NewFromInt(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountId := account1.ID
		toAccountId := account2.ID
		if i%2 == 0 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: int64(fromAccountId),
				ToAccountID:   int64(toAccountId),
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updateAccount1, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccountById(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updateAccount1.Balance, updateAccount2.Balance)

	require.True(t, account1.Balance.Equal(updateAccount1.Balance))
	require.True(t, account2.Balance.Equal(updateAccount2.Balance))

}
