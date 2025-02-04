package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Health godoc
//
//	@Summary	healthcheck
//	@Schemes	http https
//	@Tags		health
//	@Accept		json
//	@Produce	json
//	@Success	200
//	@Router		/health [get]
func (h *handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}
