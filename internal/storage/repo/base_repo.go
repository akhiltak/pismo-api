package repo

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

type baseRepo[T any] struct {
	db bun.IDB
}

func newBaseRepo[T any](db bun.IDB) *baseRepo[T] {
	return &baseRepo[T]{
		db: db,
	}
}

func (in *baseRepo[T]) Insert(ctx context.Context, model *T) (*T, error) {
	if _, err := in.db.NewInsert().Model(model).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return model, nil
}

func (in *baseRepo[T]) Update(ctx context.Context, model *T) error {
	if _, err := in.db.NewUpdate().Model(model).WherePK().Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (in *baseRepo[T]) FindByID(ctx context.Context, id int64, relation string) (*T, error) {
	model := new(T)
	query := in.db.NewSelect().Model(model)
	if relation != "" {
		query = query.Relation(relation)
	}
	if err := query.Where("id = ?", id).Scan(ctx); err != nil {
		// CHECK for the error type if it's not found and return 404.
		return nil, err
	}
	return model, nil
}

func (in *baseRepo[T]) FindByColumn(ctx context.Context, id int64, filterColumnName, relation string) ([]*T, error) {
	var models []*T
	query := in.db.NewSelect().Model(&models)
	if relation != "" {
		query = query.Relation(relation)
	}
	if err := query.Where("? = ?", bun.Ident(filterColumnName), id).Scan(ctx); err != nil {
		return nil, err
	}
	return models, nil
}

func (in *baseRepo[T]) Delete(ctx context.Context, id int64) error {
	model, err := in.FindByID(ctx, id, "") // Fetch the object first
	if err != nil {
		// CHECK for the error type if it's not found and return 404.
		return err
	}
	if _, err := in.db.NewDelete().Model(model).WherePK().Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (in *baseRepo[T]) GetAll(ctx context.Context, relation string) ([]*T, error) {
	var models []*T
	query := in.db.NewSelect().Model(&models)

	if relation != "" {
		query = query.Relation(relation)
	}

	if err := query.OrderExpr("updated_at ASC").Scan(ctx); err != nil {
		return nil, err
	}
	return models, nil
}

func (in *baseRepo[T]) RunInTx(ctx context.Context, opts *sql.TxOptions, f func(ctx context.Context, tx bun.Tx) error) error {
	return in.db.RunInTx(ctx, opts, func(ctx context.Context, tx bun.Tx) error {
		return f(ctx, tx)
	})
}
