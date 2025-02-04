package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

// Account represents customer account.
type Account struct {
	bun.BaseModel `bun:"table:accounts" swaggerignore:"true"` // Specifies the table name

	ID        int64     `json:"id" bun:"id,pk,autoincrement,type:int"`                                          // Primary key
	DocNum    string    `json:"document_number" bun:"document_number,type:varchar(255)"`                        // Document number
	CreatedAt time.Time `json:"created_at" bun:"created_at,type:timestamptz,notnull,default:current_timestamp"` // CreatedAt with default
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,type:timestamptz,notnull,default:current_timestamp"` // UpdatedAt with default
} // @name Account

var _ bun.BeforeAppendModelHook = (*Account)(nil)

func (m *Account) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now().UTC()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now().UTC()
	}
	return nil
}
