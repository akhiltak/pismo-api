package repo

import (
	"context"

	"github.com/akhiltak/pismo-api/internal/storage/models"
	"github.com/uptrace/bun"
)

type Transaction interface {
	Create(context.Context, *models.Transaction) (*models.Transaction, error)
	GetAllTransactions(context.Context) ([]*models.Transaction, error)
}

type transaction struct {
	*baseRepo[models.Transaction]
}

func NewTransactionRepo(db bun.IDB) Transaction {
	return &transaction{baseRepo: newBaseRepo[models.Transaction](db)}
}

func (a *transaction) Create(ctx context.Context, model *models.Transaction) (*models.Transaction, error) {
	return a.baseRepo.Insert(ctx, model)
}

// GetAllTransactions fetches all customer Transactions
func (a *transaction) GetAllTransactions(ctx context.Context) ([]*models.Transaction, error) {
	return a.baseRepo.GetAll(ctx, "")
}
