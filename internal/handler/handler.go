package handler

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	services "github.com/akhiltak/pismo-api/internal/service"
	"github.com/akhiltak/pismo-api/pkg/api"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Health(c echo.Context) error
	CreateAccount(c echo.Context) error
	CreateTransaction(c echo.Context) error
	GetAccountByID(c echo.Context) error
	GetTransactions(c echo.Context) error
}

type handler struct {
	transactionService services.TransactionService
}

var _ Handler = (*handler)(nil)

func New(
	transactionService services.TransactionService,
) Handler {
	return &handler{transactionService: transactionService}
}

func (h *handler) bindAndValidate(c echo.Context, obj any) error {
	ctx := h.ctx(c)
	slog.DebugContext(ctx, "binding request...")
	if err := c.Bind(obj); err != nil {
		return api.BadRequestErr("invalid request, please verify", err)
	}
	slog.DebugContext(ctx, "validating request...")
	validate := validator.New()
	if err := validate.Struct(obj); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errorMessages []string
			for _, err := range validationErrors {
				var msg string
				switch err.Tag() {
				// insert cases here when custom validations and tags are added.
				default:
					msg = fmt.Sprintf("Field validation for '%s:%s' failed on the '%s' tag", err.Field(), err.Value(), err.Tag())
				}
				errorMessages = append(errorMessages, msg)
			}
			return api.BadRequestErr(strings.Join(errorMessages, "\n"), nil)
		}
		return api.BadRequestErr(err.Error(), nil)
	}
	return nil
}

func (h *handler) ctx(c echo.Context) context.Context {
	return c.Request().Context()
}
