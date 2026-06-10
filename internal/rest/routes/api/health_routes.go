package api_routes

import (
	"database/sql"

	"github.com/Dufyz/scd-server/internal/rest/controllers"
	"github.com/labstack/echo/v4"
)

func HealthRoutes(api *echo.Group, db *sql.DB) {
	healthController := controllers.NewHealthController(db)

	api.GET("/health", healthController.GETHealth)
	api.GET("/health-status", healthController.GETHealthStatus)
}
