package adapter

import (
	"context"
	"database/sql"
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

	return &transaction, nil
}
