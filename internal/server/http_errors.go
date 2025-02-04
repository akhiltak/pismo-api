package server

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/akhiltak/pismo-api/pkg/api"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun/driver/pgdriver"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	// ctx := c.Request().Context()

	var code int
	var message string
	var internal error

	switch v := err.(type) {
	case *echo.HTTPError:
		if errors.Is(err, sql.ErrNoRows) {
			code = http.StatusNotFound
			message = "Record not found in DB"
		} else if pgErr, ok := v.Internal.(pgdriver.Error); ok {
			switch pgErr.Field('C') {
			case "23503": // Foreign key violation
				code = http.StatusBadRequest
				message = "Foreign key constraint failed: "
			case "23505": // Unique constraint violation
				code = http.StatusConflict
				message = "Duplicate record error: "
			case "22001": // String data right truncation
				code = http.StatusBadRequest
				message = "Data too long for column: "
			case "23502": // NOT NULL violation
				code = http.StatusBadRequest
				message = "Missing required field: "
			default:
				code = http.StatusInternalServerError
				message = "Database error: "
			}
			message += pgErr.Field('M')
		} else {
			code = v.Code
			message = v.Message.(string)
		}
		internal = v.Internal
	default:
		code = http.StatusInternalServerError
		message = api.InternalServerErr
		internal = err
	}
	if message == "" {
		message = api.InternalServerErr
	}

	// Return the error response in JSON format
	errorResponse := api.Response{
		Success: false,
		Code:    code,
		Error: &echo.HTTPError{
			Code:     code,
			Message:  message,
			Internal: internal,
		},
	}

	// Check if the response has already been committed
	// If not, commit the response with the error details
	if !c.Response().Committed {
		c.JSON(code, errorResponse)
	}
}
