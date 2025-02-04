package service

import (
	"context"
	"log/slog"

	"github.com/akhiltak/pismo-api/internal/storage/models"
	"github.com/akhiltak/pismo-api/internal/storage/repo"
	"github.com/akhiltak/pismo-api/pkg/api"
)

type TransactionService interface {
	CreateAccount(context.Context, *api.CreateAccountRequest) (*models.Account, error)
	GetAccountByID(context.Context, int64) (*models.Account, error)
	CreateTransaction(context.Context, *api.CreateTransactionRequest) (*models.Transaction, error)
	GetTransactions(context.Context) ([]*models.Transaction, error)
}

type txnSrv struct {
	accountRepo     repo.Account
	transactionRepo repo.Transaction
	operationRepo   repo.Operation
}

var _ TransactionService = (*txnSrv)(nil)

func NewTransactionService(
	accountRepo repo.Account,
	transactionRepo repo.Transaction,
	operationRepo repo.Operation,
) TransactionService {
	return &txnSrv{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		operationRepo:   operationRepo,
	}
}

// GetAccountByID fetches a customer account by its ID
func (s *txnSrv) GetAccountByID(ctx context.Context, id int64) (*models.Account, error) {
	return s.accountRepo.GetByID(ctx, id, true)
}

func (s *txnSrv) CreateAccount(ctx context.Context, req *api.CreateAccountRequest) (*models.Account, error) {
	return s.accountRepo.Create(ctx, &models.Account{
		DocNum: req.DocNum,
	})
}

// CreateTransaction creates a new transaction
// Truncates the amount to two decimal places before storing
// Validates the operation type exists and also finds out negative/positive amount based on credit/debit entryType
// No need to validate account ID since foreign key constraint will complain on DB insert (same for operation type id actually)
func (s *txnSrv) CreateTransaction(ctx context.Context, req *api.CreateTransactionRequest) (*models.Transaction, error) {
	// Get operation by id
	operation, err := s.operationRepo.GetByID(ctx, req.OperationTypeID, false)
	if err != nil {
		return nil, err
	}

	// Check if operation is valid
	if operation == nil {
		return nil, api.BadRequestErr(api.ErrOpTypeNotFound, nil)
	}
	// positive amount for credit and negative for debit
	switch operation.EntryType {
	case models.DebitEntry:
		req.Amount = req.Amount.Abs().Neg()
	case models.CreditEntry:
		req.Amount = req.Amount.Abs()
	}
	slog.Debug("CreateTransaction", "amount", req.Amount, "operation", operation.EntryType)

	return s.transactionRepo.Create(ctx, &models.Transaction{
		AccountID:       req.AccountID,
		OperationTypeID: req.OperationTypeID,
		Status:          models.TxnStatusCompleted,
		Amount:          req.Amount,
	})
}

func (s *txnSrv) GetTransactions(ctx context.Context) ([]*models.Transaction, error) {
	return s.transactionRepo.GetAllTransactions(ctx)
}
