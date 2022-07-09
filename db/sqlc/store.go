package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	Querier // this is due to emit_interface set to true
}

//Store provides all functions to execute db queries and transaction
//in order to use mockdb, we have to create an interface and this struct will
//satisfy that interface
type SimpleStore struct {
	*Queries // embed all the func that Queries have
	db       *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SimpleStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database Transaction
func (store *SimpleStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// tx also satisfy the DBTX interface
	q := New(tx)

	// we now have a query
	err = fn(q)

	// rollback if error
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

//TransferTxParams contains all necessary input parameters for the transfer transcation
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

//TransferTxResult contains the result of transfer transaction
type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	// these are after the balance are updated
	FromAccount Account `json:"from_account_id"`
	ToAccount   Account `json:"to_account_id"`

	//the entry records money moving out and moving in
	FromEntry Entry `json:"from_entry"`
	ToEntry   Entry `json:"to_entry"`
}

var txKey = struct{}{}

// transferTX performs a money transfer from one account to another account
// it will create a new transfer record
// add account entries, and update account balance within a single db transaction
func (store *SimpleStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(txName, "create transfer")
		//create a transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}
		// add two account entries, one for the from account, one for the to account
		// from account amount is negative
		// to account amount is positive
		fmt.Println(txName, "create entry1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// incorrect way is to get balance from db then update it (without locking)
		// the actual reason for locking is due to multiple go routines updates account in different orders
		// we should check the account ID and update them in sequence

		if arg.FromAccountID < arg.ToAccountID {
			fmt.Println(txName, "updating from account first")
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			fmt.Println(txName, "updating to account first")
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)

		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddUpdateAccountBalance(ctx, AddUpdateAccountBalanceParams{
		Amount: amount1,
		ID:     accountID1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddUpdateAccountBalance(ctx, AddUpdateAccountBalanceParams{
		Amount: amount2,
		ID:     accountID2,
	})
	if err != nil {
		return
	}

	return
}
