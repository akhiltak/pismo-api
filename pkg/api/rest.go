package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	ErrParsingID           string = "cannot parse ID, should be integer"
	ErrValidationStructure string = "cannot validate structure"
	ErrNotFound            string = "requested record not found"
	ErrOpTypeNotFound      string = "operation type record not found"
	InternalServerErr      string = "Somewhere something went wrong but don't worry, we are on it."
)

type Response struct {
	Success bool            `json:"success"`
	Code    int             `json:"code,omitempty"`
	Error   *echo.HTTPError `json:"error,omitempty"`
	// Data    any             `json:"data,omitempty"`
} // @name Response

func CustomErr(code int, msg string, err error) *echo.HTTPError {
	return &echo.HTTPError{
		Code:     code,
		Message:  msg,
		Internal: err,
	}
}

func BadRequestErr(msg string, err error) *echo.HTTPError {
	return CustomErr(http.StatusBadRequest, msg, err)
}

func UnauthorizedErr(msg string, err error) *echo.HTTPError {
	return CustomErr(http.StatusUnauthorized, msg, err)
}

func ForbiddenErr(msg string, err error) *echo.HTTPError {
	return CustomErr(http.StatusForbidden, msg, err)
}

func NotFoundErr(msg string, err error) *echo.HTTPError {
	return CustomErr(http.StatusNotFound, msg, err)
}

func ServerErr(err error) *echo.HTTPError {
	return CustomErr(http.StatusInternalServerError, InternalServerErr, err)
}

// JWT data keys
const JWTBodyRequestContextKey string = "jwtBody"

// error messages
// here var used in place of const to allow for capitalized error message
var ErrMsgInvalidJWT string = "Invalid JWT given"
var ErrMsgMalformedData string = "Malformed data given as input"

// application global errors
var ErrMalformedData error = errors.New(ErrMsgMalformedData)
