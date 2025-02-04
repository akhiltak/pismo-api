package handler

import (
	"log/slog"
	"net/http"

	_ "github.com/akhiltak/pismo-api/internal/storage/models"
	"github.com/akhiltak/pismo-api/pkg/api"
	"github.com/labstack/echo/v4"
)

// CreateTransaction godoc
//
//	@Summary	CreateTransaction
//	@Schemes	http https
//	@Tags		transaction
//	@Accept		json
//	@Produce	json
//	@Param		request	body		api.CreateTransactionRequest	true	"CreateTransactionRequest"
//	@Success	201		{object}	models.Transaction
//	@Failure	400		{object}	api.Response
//	@Failure	500		{object}	api.Response
//	@Router		/transactions [post]
func (h *handler) CreateTransaction(c echo.Context) error {
	req := &api.CreateTransactionRequest{}
	if err := h.bindAndValidate(c, req); err != nil {
		return err
	}
	slog.Debug("CreateTransaction", "req", *req)

	transaction, err := h.transactionService.CreateTransaction(c.Request().Context(), req)
	if err != nil {
		return api.ServerErr(err)
	}
	return c.JSON(http.StatusCreated, transaction)
}

// GetTransactions godoc
//
//	@Summary	GetTransactions
//	@Schemes	http https
//	@Tags		transaction
//	@Accept		json
//	@Produce	json
//	@Success	200	{array}		models.Transaction
//	@Failure	400	{object}	api.Response
//	@Failure	500	{object}	api.Response
//	@Router		/transactions [get]
func (h *handler) GetTransactions(c echo.Context) error {
	slog.Debug("GetTransactions")
	transactions, err := h.transactionService.GetTransactions(c.Request().Context())
	if err != nil {
		return api.ServerErr(err)
	}
	return c.JSON(http.StatusOK, transactions)
}
