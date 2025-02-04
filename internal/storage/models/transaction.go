package models

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type TxnStatus string // @name TxnStatus

const (
	TxnStatusPending   TxnStatus = "pending"
	TxnStatusCompleted TxnStatus = "completed"
	TxnStatusFailed    TxnStatus = "failed"
)

func (ts TxnStatus) String() string {
	return string(ts)
}

func (ts TxnStatus) Validate() error {
	switch ts {
	case TxnStatusPending, TxnStatusCompleted, TxnStatusFailed:
		return nil
	default:
		return fmt.Errorf("invalid status: %s", ts)
	}
}

// Transaction represents customer transactions.
type Transaction struct {
	bun.BaseModel `bun:"table:transactions" swaggerignore:"true"` // Specifies the table name

	ID              int64           `json:"id" bun:"id,pk,autoincrement,type:int"`                                          // Primary key
	AccountID       int64           `json:"account_id" bun:"account_id,type:int,notnull"`                                   // Foreign key to account
	OperationTypeID int64           `json:"operationTypeID" bun:"operation_type_id,type:int,notnull"`                       // Foreign key to OperationType
	Amount          decimal.Decimal `json:"amount" bun:"amount,type:float8,notnull"`                                        // Transaction amount
	Status          TxnStatus       `json:"status" bun:"status,type:varchar(255),notnull"`                                  // status
	EventDate       time.Time       `json:"event_date" bun:"event_date,type:timestamptz,notnull,default:current_timestamp"` // CreatedAt with default, called EventDate due to assignment instructions
	UpdatedAt       time.Time       `json:"updated_at" bun:"updated_at,type:timestamptz,notnull,default:current_timestamp"` // UpdatedAt with default
} // @name Transaction

var _ bun.BeforeAppendModelHook = (*Transaction)(nil)

func (m *Transaction) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.EventDate = time.Now().UTC()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now().UTC()
	}
	return nil
}
