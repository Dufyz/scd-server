package api_routes

import (
	"github.com/Dufyz/scd-server/internal/rest/controllers"
	"github.com/labstack/echo/v4"
)

func HealthRoutes(api *echo.Group) {
	healthController := controllers.NewHealthController()

	api.GET("/health", healthController.GETHealth)
	api.GET("/health/kafka", healthController.GETKafkaHealth)
}
