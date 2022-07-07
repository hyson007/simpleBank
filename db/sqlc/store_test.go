package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// run n concurrent transfer transaction
	// as we want to make sure it works in concurrent
	// test with 5 go routine, each transfer 10
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func(txName string) {
			tp := TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			}
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, tp)
			// if err != nil {
			errs <- err
			// } else {
			results <- result
			// }
		}(txName)
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		// check each result
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		//check if transfer is really created in DB
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries from the result
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, amount, -fromEntry.Amount)
		require.NotZero(t, fromEntry.CreatedAt)

		//check entries from db
		// fmt.Println(fromEntry)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		//check entries from the result
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.CreatedAt)

		//check entries from db
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		//check account balance
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}

	// check the updated balance
	updatedAccount1, err := testQuery.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQuery.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// run n concurrent transfer transaction
	// as we want to make sure it works in concurrent
	// test with 5 go routine, each transfer 10
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		FromAccountID := account1.ID
		ToAccountID := account2.ID

		//make half of the transaction to be from account2 to account1
		if i%2 == 1 {
			FromAccountID = account2.ID
			ToAccountID = account1.ID
		}

		// the actual reason for locking is due to multiple go routines updates account in different orders
		// we should check the account ID and update them in sequence

		go func(txName string) {
			tp := TransferTxParams{
				FromAccountID: FromAccountID,
				ToAccountID:   ToAccountID,
				Amount:        amount,
			}
			ctx := context.WithValue(context.Background(), txKey, txName)
			_, err := store.TransferTx(ctx, tp)
			// if err != nil {
			errs <- err
			// } else {
			// }
		}(txName)
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the updated balance
	updatedAccount1, err := testQuery.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQuery.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, updatedAccount1.Balance, account1.Balance)
	require.Equal(t, updatedAccount2.Balance, account2.Balance)
}
