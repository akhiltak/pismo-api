package repo

import (
	"context"

	"github.com/akhiltak/pismo-api/internal/storage/models"
	"github.com/uptrace/bun"
)

type Account interface {
	Create(context.Context, *models.Account) (*models.Account, error)
	GetAllAccounts(context.Context) ([]*models.Account, error)
	GetByID(context.Context, int64, bool) (*models.Account, error)
}

type account struct {
	*baseRepo[models.Account]
}

func NewAccountRepo(db bun.IDB) Account {
	return &account{baseRepo: newBaseRepo[models.Account](db)}
}

func (a *account) Create(ctx context.Context, model *models.Account) (*models.Account, error) {
	return a.baseRepo.Insert(ctx, model)
}

// GetAllAccounts fetches all customer Accounts
func (a *account) GetAllAccounts(ctx context.Context) ([]*models.Account, error) {
	return a.baseRepo.GetAll(ctx, "")
}

// GetByID fetches an Account by ID
func (a *account) GetByID(ctx context.Context, id int64, associations bool) (*models.Account, error) {
	return a.baseRepo.FindByID(ctx, id, "")
}
