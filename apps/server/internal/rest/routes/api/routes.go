package api_routes

import (
	db "github.com/Dufyz/scd-server/infra/database"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

func SetupApiRoutes(api *echo.Group, conn *db.ReplicatedDB, queueClient *asynq.Client) {
	ChatRoutes(api, conn)
	MessageRoutes(api, conn)
	HealthRoutes(api)
}
