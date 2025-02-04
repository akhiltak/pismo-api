package api

import "github.com/shopspring/decimal"

type CreateAccountRequest struct {
	DocNum string `json:"document_number" validate:"required"`
} // @name CreateAccountRequest

type CreateTransactionRequest struct {
	AccountID       int64           `json:"account_id" validate:"required"`
	OperationTypeID int64           `json:"operation_type_id" validate:"required"`
	Amount          decimal.Decimal `json:"amount" validate:"required"`
} // @name CreateTransactionRequest
