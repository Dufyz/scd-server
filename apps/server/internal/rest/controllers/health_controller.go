package controllers

import (
	"net"
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
	broker := brokers[0]

	conn, err := net.DialTimeout("tcp", broker, 5*time.Second)
	if err != nil {
		return ctx.JSON(http.StatusServiceUnavailable, map[string]string{
			"status":  "down",
			"brokers": broker,
			"error":   err.Error(),
		})
	}
	conn.Close()

	return ctx.JSON(http.StatusOK, map[string]string{
		"status":  "up",
		"brokers": broker,
	})
}
