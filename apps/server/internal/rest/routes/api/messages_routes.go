package api_routes

import (
	"database/sql"

	"github.com/Dufyz/scd-server/infra/database/repositories"
	"github.com/Dufyz/scd-server/internal/domain/usecases"
	"github.com/Dufyz/scd-server/internal/rest/controllers"
	"github.com/labstack/echo/v4"
)

func MessageRoutes(api *echo.Group, db *sql.DB) {
	messageRepository := repositories.NewMessageRepository(db)

	messageUsecase := usecases.NewMessageUsecase(messageRepository)

	messageController := controllers.NewMessageController(messageUsecase)

	api.POST("/messages", messageController.POSTMessage)

	api.PUT("/messages/:id", messageController.PUTMessage)

	api.DELETE("/messages/:id", messageController.DELETEMessage)
}
