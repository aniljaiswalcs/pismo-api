package adapter

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/aniljaiswalcs/pismo/model"
)

type TransactionRepositoryPostgres struct {
	db *sql.DB
}

func NewTransactionRepositoryPostgres(db *sql.DB) *TransactionRepositoryPostgres {
	return &TransactionRepositoryPostgres{
		db: db,
	}
}

func (t *TransactionRepositoryPostgres) CreateTransaction(ctx context.Context, transaction model.Transaction) (*model.Transaction, error) {

	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := "INSERT INTO transactions (account_id, operation_type_id, amount) VALUES ($1, $2, $3) RETURNING transaction_id"
	err := t.db.QueryRowContext(
		ctxTimeout,
		query,
		transaction.AccountId,
		transaction.OperationTypeId,
		transaction.Amount).
		Scan(&transaction.TransactionId)

	if err != nil {
		log.Printf("TransactionRepositoryPostgres#CreateTransaction: Database query (%s) failed: %s", query, err)
		return nil, err
	}
	//update values
	if transaction.OperationTypeId == 4 {
		err := t.SubtractTransaction(ctx, transaction)
		if err != nil {
			return nil, err
		}
	}
	// find the updated transaction valuues
	transactionId, err := t.FindAccount(ctx, transaction.TransactionId)
	if err != nil {
		return nil, err
	}
	return transactionId, nil
}

func (t *TransactionRepositoryPostgres) SubtractTransaction(ctx context.Context, transaction model.Transaction) error {

	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// fetch amount using account id sort by time;

	rows, err := t.db.Query("SELECT transaction_id, balance, account_id, operation_type_id FROM transactions WHERE account_id < 4 sort by EventDate DESC")
	if err != nil {
		fmt.Println(" Error during SubtractTransaction query")
		return nil
	}

	result := []model.Transaction{}

	defer rows.Close()
	for rows.Next() {
		res := model.Transaction{} // creating new struct for every row
		err = rows.Scan(&res.TransactionId, &res.Balance, &res.AccountId, &res.OperationTypeId)
		if err != nil {
			log.Println(err)
		}
		result = append(result, res) // add new row information
	}

	initialVal := float64(transaction.Amount)

	for index, _ := range result {
		balance := math.Abs(float64(result[index].Balance))
		if initialVal-balance > 0 {
			initialVal -= balance
			result[index].Balance = 0
		} else {
			break
		}
	}
	transaction.Balance = float32(initialVal)
	err = t.UpdateTransactiondatabse(ctxTimeout, result, transaction)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return nil
}

func (t *TransactionRepositoryPostgres) UpdateTransactiondatabse(ctx context.Context, result []model.Transaction, initialtransaction model.Transaction) error {

	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := "update transactions balance VALUES $1 where transaction_id, account_id, operation_type_id VALUES ($2, $3, $4)"

	for index, _ := range result {
		err := t.db.QueryRowContext(
			ctxTimeout,
			query,
			result[index].Balance,
			result[index].TransactionId,
			result[index].AccountId,
			result[index].OperationTypeId)

		if err != nil {
			log.Printf("TransactionRepositoryPostgres#CreateTransaction: Database query (%s) failed: %s", query, err)
			return nil
		}
	}

	// update initial balance
	err := t.db.QueryRowContext(
		ctxTimeout,
		query,
		initialtransaction.Balance,
		initialtransaction.TransactionId,
		initialtransaction.AccountId,
		initialtransaction.OperationTypeId)

	if err != nil {
		log.Printf("TransactionRepositoryPostgres#CreateTransaction: Database query (%s) failed: %s", query, err)
		return nil
	}

	return nil
}

func (t *TransactionRepositoryPostgres) FindAccount(ctx context.Context, accountId uint64) (*model.Transaction, error) {

	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	transaction := model.Transaction{}
	query := "SELECT account_id, operation_type_id, amount,balance FROM transactions WHERE transaction_id=$1 LIMIT 1"
	result := t.db.QueryRowContext(ctxTimeout, query, accountId)
	err := result.Scan(&transaction.AccountId, &transaction.Amount, &transaction.Balance, &transaction.OperationTypeId, &transaction.TransactionId)
	if err != nil {
		log.Printf("transactionRepositoryPostgres#FindAccount: Database query (%s) failed: %s", query, err)

		return nil, err
	}

	return &transaction, nil
}
