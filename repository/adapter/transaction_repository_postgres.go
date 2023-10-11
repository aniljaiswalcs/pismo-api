package adapter

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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

	query := "INSERT INTO transactions (account_id, operation_type_id, amount, balance) VALUES ($1, $2, $3, $3) RETURNING transaction_id"
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
	transactionId, err := t.FindtransactionAccount(ctx, transaction.TransactionId)
	if err != nil {
		return nil, err
	}
	return transactionId, nil
}

func (t *TransactionRepositoryPostgres) SubtractTransaction(ctx context.Context, transaction model.Transaction) error {

	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// fetch amount using account id sort by time;
	query := "SELECT transaction_id, balance, account_id, operation_type_id FROM transactions WHERE account_id = $1 AND operation_type_id < 4 order by created_at DESC"

	rows, err := t.db.Query(query, transaction.AccountId)
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

	initialVal := transaction.Amount
	for index, _ := range result {
		res := result[index].Balance + initialVal
		if res > 0 {
			result[index].Balance = 0
			initialVal = res
		} else if res <= 0 {
			initialVal = 0
			break
		}
	}
	transaction.Balance = initialVal
	err = t.UpdateTransactiondatabse(ctxTimeout, result, transaction)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return nil
}

func (t *TransactionRepositoryPostgres) UpdateTransactiondatabse(ctx context.Context, result []model.Transaction, initialtransaction model.Transaction) error {

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "UPDATE transactions set balance = $1 where transaction_id = $2 AND account_id = $3 AND operation_type_id = $4"

	for _, res := range result {
		_, err := t.db.ExecContext(
			ctxTimeout,
			query,
			res.Balance,
			res.TransactionId,
			res.AccountId,
			res.OperationTypeId)

		if err != nil {
			log.Printf("TransactionRepositoryPostgres#UpdateTransaction: Database query (%s) failed: %s", query, err)
			return nil
		}
	}

	// update initial balance
	_, err := t.db.ExecContext(
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

func (t *TransactionRepositoryPostgres) FindtransactionAccount(ctx context.Context, transactionid uint64) (*model.Transaction, error) {

	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	transaction := model.Transaction{}
	query := "SELECT account_id, operation_type_id, amount,balance, transaction_id FROM transactions WHERE transaction_id=$1 LIMIT 1"
	result := t.db.QueryRowContext(ctxTimeout, query, transactionid)
	err := result.Scan(&transaction.AccountId, &transaction.OperationTypeId, &transaction.Amount, &transaction.Balance, &transaction.TransactionId)
	if err != nil {
		log.Printf("transactionRepositoryPostgres#FindAccount: Database query (%s) failed: %s", query, err)

		return nil, err
	}

	return &transaction, nil
}
