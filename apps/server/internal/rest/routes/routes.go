package routes

import (
	db "github.com/Dufyz/scd-server/infra/database"
	api_routes "github.com/Dufyz/scd-server/internal/rest/routes/api"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, conn *db.ReplicatedDB, queueClient *asynq.Client) {
	api := e.Group("/api")
	api_routes.SetupApiRoutes(api, conn, queueClient)
}
