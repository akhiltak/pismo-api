package models

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type EntryType string // @name EntryType

const (
	CreditEntry EntryType = "credit"
	DebitEntry  EntryType = "debit"
)

func (et EntryType) String() string {
	return string(et)
}

func (et EntryType) Validate() error {
	switch et {
	case CreditEntry, DebitEntry:
		return nil
	default:
		return fmt.Errorf("invalid operation type: %s", et)
	}
}

// OperationType represents type of operations.
type OperationType struct {
	bun.BaseModel `bun:"table:operation_types" swaggerignore:"true"` // Specifies the table name

	ID          int64     `json:"id" bun:"id,pk,autoincrement,type:int"`                                          // Primary key
	Description string    `json:"description" bun:"description,type:varchar(255)"`                                // Description
	EntryType   EntryType `json:"type" bun:"entry_type,type:varchar(255)"`                                        // type (credit/debit)
	CreatedAt   time.Time `json:"created_at" bun:"created_at,type:timestamptz,notnull,default:current_timestamp"` // CreatedAt with default
	UpdatedAt   time.Time `json:"updated_at" bun:"updated_at,type:timestamptz,notnull,default:current_timestamp"` // UpdatedAt with default
} // @name OperationType

var _ bun.BeforeAppendModelHook = (*OperationType)(nil)

func (m *OperationType) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now().UTC()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now().UTC()
	}
	return nil
}
