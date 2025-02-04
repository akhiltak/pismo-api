package repo

import (
	"context"

	"github.com/akhiltak/pismo-api/internal/storage/models"
	"github.com/uptrace/bun"
)

type Operation interface {
	GetByID(context.Context, int64, bool) (*models.OperationType, error)
}

type operation struct {
	*baseRepo[models.OperationType]
}

func NewOperationRepo(db bun.IDB) Operation {
	return &operation{baseRepo: newBaseRepo[models.OperationType](db)}
}

// GetByID fetches an Operation by ID
func (o *operation) GetByID(ctx context.Context, id int64, associations bool) (*models.OperationType, error) {
	return o.baseRepo.FindByID(ctx, id, "")
}
