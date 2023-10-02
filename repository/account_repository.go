package repository

import (
	"context"

	"github.com/aniljaiswalcs/pismo/model"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, account model.Account) (*model.Account, error)
	FindAccount(ctx context.Context, accountId uint64) (*model.Account, error)
}
