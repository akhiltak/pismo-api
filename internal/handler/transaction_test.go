package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mockService "github.com/akhiltak/pismo-api/internal/service/mock_services"
	"github.com/akhiltak/pismo-api/internal/storage/models"
	"github.com/akhiltak/pismo-api/pkg/api"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockService.NewMockTransactionService(ctrl)
	h := &handler{transactionService: mockService}

	e := echo.New()

	t.Run("successful creation", func(t *testing.T) {
		reqBody := `{"account_id":1,"operation_type_id":1,"amount":100.50}`
		req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(&models.Transaction{
			ID:              1,
			AccountID:       1,
			OperationTypeID: 1,
			Amount:          decimal.NewFromFloat(100.50),
			Status:          models.TxnStatusCompleted,
		}, nil)

		if assert.NoError(t, h.CreateTransaction(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			var response models.Transaction
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, int64(1), response.ID)
		}
	})

	t.Run("invalid request", func(t *testing.T) {
		reqBody := `{"account_id":"invalid","operation_type_id":1,"amount":100.50}`
		req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateTransaction(c)
		assert.Error(t, err)
		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, he.Code)
	})

	t.Run("invalid request operation type", func(t *testing.T) {
		reqBody := `{"account_id":1,"operation_type_id":0,"amount":100.50}`
		req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateTransaction(c)
		assert.Error(t, err)
		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, he.Code)
	})
}

func TestGetTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockService.NewMockTransactionService(ctrl)
	h := &handler{transactionService: mockService}

	e := echo.New()

	t.Run("successful retrieval", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockTransactions := []*models.Transaction{
			{ID: 1, AccountID: 1, OperationTypeID: 1, Amount: decimal.NewFromFloat(100.50), Status: models.TxnStatusCompleted},
			{ID: 2, AccountID: 2, OperationTypeID: 2, Amount: decimal.NewFromFloat(200.75), Status: models.TxnStatusPending},
		}

		mockService.EXPECT().GetTransactions(gomock.Any()).Return(mockTransactions, nil)

		if assert.NoError(t, h.GetTransactions(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var response []*models.Transaction
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response, 2)
			assert.Equal(t, int64(1), response[0].ID)
			assert.Equal(t, int64(2), response[1].ID)
		}
	})

	t.Run("service error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().GetTransactions(gomock.Any()).Return(nil, api.ServerErr(nil))

		err := h.GetTransactions(c)
		assert.Error(t, err)
		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, he.Code)
	})
}
