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

type result struct {
	AccountId       uint64  `json:"account_id"`
	OperationTypeId uint32  `json:"operation_type_id"`
	Balance         float32 `json:"balance"`
}

func NewTransactionRepositoryPostgres(db *sql.DB) *TransactionRepositoryPostgres {
	return &TransactionRepositoryPostgres{
		db: db,
	}
}

func (t *TransactionRepositoryPostgres) CreateTransaction(ctx context.Context, transaction model.Transaction) (*model.Transaction, error) {

	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if transaction.OperationTypeId == 4 {
		t.SubtractTransaction(ctx, transaction)
	}
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

func (t *TransactionRepositoryPostgres) SubtractTransaction(ctx context.Context, transaction model.Transaction) {
	// fetch amount using account id sort by time;

	rows, err := t.db.Query("SELECT balance, account_id, operation_type_id FROM transactions WHERE account_id < 4 sort by EventDate DESC")
	if err != nil {
		fmt.Println(" Error during SubtractTransaction query")
		return
	}

	result := []model.Transaction{}

	defer rows.Close()
	for rows.Next() {
		res := model.Transaction{} // creating new struct for every row
		err = rows.Scan(&res.Balance, &res.AccountId, &res.OperationTypeId)
		if err != nil {
			log.Println(err)
		}
		result = append(result, res) // add new row information
	}

	initialVal := float64(transaction.Amount) //10

	for index, _ := range result {

		if initialVal-math.Abs(result[index].Balance) > 0 {
			initialVal -= result[index].Balance
			result[index].Balance = 0
		} else {
			break
		}

	}

}
