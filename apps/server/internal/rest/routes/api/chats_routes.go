package api_routes

import (
	"database/sql"

	"github.com/Dufyz/scd-server/infra/database/repositories"
	"github.com/Dufyz/scd-server/internal/domain/usecases"
	"github.com/Dufyz/scd-server/internal/rest/controllers"
	"github.com/labstack/echo/v4"
)

func ChatRoutes(api *echo.Group, db *sql.DB) {
	chatRepository := repositories.NewChatRepository(db)
	messageRepository := repositories.NewMessageRepository(db)

	chatUsecase := usecases.NewChatUsecase(chatRepository)
	messageUsecase := usecases.NewMessageUsecase(messageRepository)

	chatController := controllers.NewChatController(chatUsecase, messageUsecase)

	api.GET("/chats", chatController.GETChats)
	api.GET("/chats/:id/messages", chatController.GETChatMessages)

	api.POST("/chats", chatController.POSTChat)

	api.PUT("/chats/:id", chatController.PUTChat)

	api.DELETE("/chats/:id", chatController.DELETEChat)
}
