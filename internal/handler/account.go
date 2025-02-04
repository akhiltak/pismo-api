package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	_ "github.com/akhiltak/pismo-api/internal/storage/models"
	"github.com/akhiltak/pismo-api/pkg/api"
	"github.com/labstack/echo/v4"
)

// CreateAccount godoc
//
//	@Summary	CreateAccount
//	@Schemes	http https
//	@Tags		account
//	@Accept		json
//	@Produce	json
//	@Param		request	body		api.CreateAccountRequest	true	"CreateAccountRequest"
//	@Success	201		{object}	models.Account
//	@Failure	400		{object}	api.Response
//	@Failure	500		{object}	api.Response
//	@Router		/accounts [post]
func (h *handler) CreateAccount(c echo.Context) error {
	req := &api.CreateAccountRequest{}
	if err := h.bindAndValidate(c, req); err != nil {
		return err
	}
	slog.Debug("CreateAccount", "request", *req)

	account, err := h.transactionService.CreateAccount(c.Request().Context(), req)
	if err != nil {
		return api.ServerErr(err)
	}
	return c.JSON(http.StatusCreated, account)
}

// GetAccountByID godoc
//
//	@Summary	GetAccountByID
//	@Schemes	http https
//	@Tags		account
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Account ID"
//	@Success	200	{object}	models.Account
//	@Failure	400	{object}	api.Response
//	@Failure	404	{object}	api.Response
//	@Failure	500	{object}	api.Response
//	@Router		/accounts/{id} [get]
func (h *handler) GetAccountByID(c echo.Context) error {
	idStr := c.Param("id")
	slog.Debug("GetAccountByID", "id", idStr)

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return api.BadRequestErr(api.ErrParsingID, err)
	}

	account, err := h.transactionService.GetAccountByID(c.Request().Context(), id)
	if err != nil {
		return api.ServerErr(err)
	}
	return c.JSON(http.StatusOK, account)
}
