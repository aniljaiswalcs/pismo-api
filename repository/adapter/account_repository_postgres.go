package adapter

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/aniljaiswalcs/pismo/model"
)

type AccountRepositoryPostgres struct {
	db *sql.DB
}

func NewAccountRepositoryPostgres(db *sql.DB) *AccountRepositoryPostgres {
	return &AccountRepositoryPostgres{
		db: db,
	}
}

func (a *AccountRepositoryPostgres) CreateAccount(ctx context.Context, account model.Account) (*model.Account, error) {

	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := "INSERT INTO accounts (document_number) VALUES ($1) RETURNING account_id"

	err := a.db.QueryRowContext(ctxTimeout, query, account.DocumentNumber).Scan(&account.AccountId)
	if err != nil {
		log.Printf("AccountRepositoryPostgres#CreateAccount: Database query (%s) failed: %s", query, err)
		return nil, err
	}

	return &account, nil

}

func (a *AccountRepositoryPostgres) FindAccount(ctx context.Context, accountId uint64) (*model.Account, error) {

	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	account := model.Account{}
	query := "SELECT account_id, document_number FROM accounts WHERE account_id=$1 LIMIT 1"
	result := a.db.QueryRowContext(ctxTimeout, query, accountId)
	err := result.Scan(&account.AccountId, &account.DocumentNumber)
	if err != nil {
		log.Printf("AccountRepositoryPostgres#FindAccount: Database query (%s) failed: %s", query, err)

		return nil, err
	}

	return &account, nil
}
