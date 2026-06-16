package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) GETHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Health status endpoint is operational",
		"status":  "up",
	})
}
