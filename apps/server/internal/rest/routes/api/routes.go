package api_routes

import (
	"database/sql"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

func SetupApiRoutes(api *echo.Group, db *sql.DB, queueClient *asynq.Client) {
	ChatRoutes(api, db)
	MessageRoutes(api, db)
	HealthRoutes(api, db)
}
