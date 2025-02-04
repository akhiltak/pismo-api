package service

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/akhiltak/pismo-api/internal/storage/models"
	mockRepo "github.com/akhiltak/pismo-api/internal/storage/repo/mock_repo"
	"github.com/akhiltak/pismo-api/pkg/api"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAccountRepo := mockRepo.NewMockAccount(ctrl)
	service := NewTransactionService(mockAccountRepo, nil, nil)

	t.Run("successful creation", func(t *testing.T) {
		req := &api.CreateAccountRequest{DocNum: "12345678"}
		expectedAccount := &models.Account{ID: 1, DocNum: "12345678"}

		mockAccountRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(expectedAccount, nil)

		account, err := service.CreateAccount(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, expectedAccount, account)
	})

	t.Run("repo error", func(t *testing.T) {
		req := &api.CreateAccountRequest{DocNum: "12345678"}

		mockAccountRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)

		account, err := service.CreateAccount(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, account)
	})
}

func TestGetAccountByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAccountRepo := mockRepo.NewMockAccount(ctrl)
	service := NewTransactionService(mockAccountRepo, nil, nil)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedAccount := &models.Account{ID: 1, DocNum: "12345678"}

		mockAccountRepo.EXPECT().GetByID(gomock.Any(), int64(1), true).Return(expectedAccount, nil)

		account, err := service.GetAccountByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedAccount, account)
	})

	t.Run("repo error", func(t *testing.T) {
		mockAccountRepo.EXPECT().GetByID(gomock.Any(), int64(1), true).Return(nil, assert.AnError)

		account, err := service.GetAccountByID(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, account)
	})
}

func TestCreateTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := mockRepo.NewMockTransaction(ctrl)
	mockOperationRepo := mockRepo.NewMockOperation(ctrl)
	service := NewTransactionService(nil, mockTransactionRepo, mockOperationRepo)

	// operation types
	op1 := &models.OperationType{ID: 1, Description: "Normal Purchase", EntryType: models.DebitEntry}
	// op2 := &models.OperationType{ID: 2, Description: "Purchase with installments", EntryType: models.DebitEntry}
	// op3 := &models.OperationType{ID: 3, Description: "Withdrawal", EntryType: models.DebitEntry}
	op4 := &models.OperationType{ID: 4, Description: "Credit Voucher", EntryType: models.CreditEntry}

	t.Run("successful creation - debit", func(t *testing.T) {
		req := &api.CreateTransactionRequest{
			AccountID:       1,
			OperationTypeID: 1,
			Amount:          decimal.NewFromFloat(100.50),
		}

		expectedTransaction := &models.Transaction{
			ID:              1,
			AccountID:       1,
			OperationTypeID: 1,
			Amount:          decimal.NewFromFloat(-100.50),
			Status:          models.TxnStatusCompleted,
		}

		mockOperationRepo.EXPECT().GetByID(gomock.Any(), int64(1), false).Return(op1, nil)
		mockTransactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(expectedTransaction, nil)

		transaction, err := service.CreateTransaction(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, expectedTransaction, transaction)
	})

	t.Run("successful creation - credit", func(t *testing.T) {
		req := &api.CreateTransactionRequest{
			AccountID:       1,
			OperationTypeID: 2,
			Amount:          decimal.NewFromFloat(100.50),
		}
		expectedTransaction := &models.Transaction{
			ID:              1,
			AccountID:       1,
			OperationTypeID: 2,
			Amount:          decimal.NewFromFloat(100.50),
			Status:          models.TxnStatusCompleted,
		}

		mockOperationRepo.EXPECT().GetByID(gomock.Any(), int64(2), false).Return(op4, nil)
		mockTransactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(expectedTransaction, nil)

		transaction, err := service.CreateTransaction(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, expectedTransaction, transaction)
	})

	t.Run("invalid operation type", func(t *testing.T) {
		req := &api.CreateTransactionRequest{
			AccountID:       1,
			OperationTypeID: 3,
			Amount:          decimal.NewFromFloat(100.50),
		}

		mockOperationRepo.EXPECT().GetByID(gomock.Any(), int64(3), false).Return(nil, fmt.Errorf("wrong, something is!!!"))

		transaction, err := service.CreateTransaction(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, transaction)
	})

	t.Run("missing operation type", func(t *testing.T) {
		req := &api.CreateTransactionRequest{
			AccountID:       1,
			OperationTypeID: 3,
			Amount:          decimal.NewFromFloat(100.50),
		}

		mockOperationRepo.EXPECT().GetByID(gomock.Any(), int64(3), false).Return(nil, nil)

		transaction, err := service.CreateTransaction(context.Background(), req)
		assert.Error(t, err)
		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, he.Code)
		assert.Nil(t, transaction)
	})
}

func TestGetTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := mockRepo.NewMockTransaction(ctrl)
	service := NewTransactionService(nil, mockTransactionRepo, nil)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedTransactions := []*models.Transaction{
			{ID: 1, AccountID: 1, Amount: decimal.NewFromFloat(100.50)},
			{ID: 2, AccountID: 2, Amount: decimal.NewFromFloat(-50.25)},
		}

		mockTransactionRepo.EXPECT().GetAllTransactions(gomock.Any()).Return(expectedTransactions, nil)

		transactions, err := service.GetTransactions(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedTransactions, transactions)
	})

	t.Run("repo error", func(t *testing.T) {
		mockTransactionRepo.EXPECT().GetAllTransactions(gomock.Any()).Return(nil, assert.AnError)

		transactions, err := service.GetTransactions(context.Background())
		assert.Error(t, err)
		assert.Nil(t, transactions)
	})
}
