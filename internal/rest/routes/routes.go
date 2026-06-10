package routes

import (
	"database/sql"

	api_routes "github.com/Dufyz/scd-server/internal/rest/routes/api"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, db *sql.DB, queueClient *asynq.Client) {
	api := e.Group("/api")
	api_routes.SetupApiRoutes(api, db, queueClient)
}
