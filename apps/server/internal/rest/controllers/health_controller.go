package controllers

import (
	"context"
	"net/http"
	"time"

	kafkaInfra "github.com/Dufyz/scd-server/infra/kafka"
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

func (c *HealthController) GETKafkaHealth(ctx echo.Context) error {
	brokers := kafkaInfra.BrokersFromEnv()
	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), 5*time.Second)
	defer cancel()

	err := kafkaInfra.Produce(reqCtx, brokers, "health-check", []byte("ping"), []byte("ping"))
	if err != nil {
		return ctx.JSON(http.StatusServiceUnavailable, map[string]string{
			"status":  "down",
			"brokers": brokers[0],
			"error":   err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"status":  "up",
		"brokers": brokers[0],
	})
}
