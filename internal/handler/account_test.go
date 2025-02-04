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
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockService.NewMockTransactionService(ctrl)
	h := &handler{transactionService: mockService}

	e := echo.New()

	t.Run("successful creation", func(t *testing.T) {
		reqBody := `{"document_number":"12345678"}`
		req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(&models.Account{
			ID:     1,
			DocNum: "12345678",
		}, nil)

		if assert.NoError(t, h.CreateAccount(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			var response models.Account
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, int64(1), response.ID)
			assert.Equal(t, "12345678", response.DocNum)
		}
	})

	t.Run("invalid request", func(t *testing.T) {
		reqBody := `{"document_number":""}`
		req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateAccount(c)
		assert.Error(t, err)
		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, he.Code)
	})

	t.Run("service error", func(t *testing.T) {
		reqBody := `{"document_number":"12345678"}`
		req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(nil, api.ServerErr(nil))

		err := h.CreateAccount(c)
		assert.Error(t, err)
		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, he.Code)
	})
}

func TestGetAccountByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockService.NewMockTransactionService(ctrl)
	h := &handler{transactionService: mockService}

	e := echo.New()

	t.Run("successful retrieval", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/accounts/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockService.EXPECT().GetAccountByID(gomock.Any(), int64(1)).Return(&models.Account{
			ID:     1,
			DocNum: "12345678",
		}, nil)

		if assert.NoError(t, h.GetAccountByID(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var response models.Account
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, int64(1), response.ID)
			assert.Equal(t, "12345678", response.DocNum)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/accounts/invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/accounts/:id")
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := h.GetAccountByID(c)
		assert.Error(t, err)
		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, he.Code) // 400 status code, bad request
	})

	t.Run("service error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/accounts/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockService.EXPECT().GetAccountByID(gomock.Any(), int64(1)).Return(nil, api.ServerErr(nil))

		err := h.GetAccountByID(c)
		assert.Error(t, err)
		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, he.Code) // 500 status code, internal server error
	})
}
