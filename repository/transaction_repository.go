package repository

import (
	"context"

	"github.com/aniljaiswalcs/pismo/model"
)

type TransactionRepository interface {
	CreateTransaction(context.Context, model.Transaction) (*model.Transaction, error)
}
