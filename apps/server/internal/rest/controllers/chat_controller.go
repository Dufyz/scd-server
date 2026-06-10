package controllers

import (
	"net/http"

	"github.com/Dufyz/scd-server/infra/log"
	"github.com/Dufyz/scd-server/internal/domain/usecases"
	"github.com/Dufyz/scd-server/internal/shared/dtos"
	"github.com/Dufyz/scd-server/internal/shared/errors"
	"github.com/Dufyz/scd-server/internal/shared/utils"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type chatController struct {
	usecase        usecases.ChatUsecase
	messageUsecase usecases.MessageUsecase
}

func NewChatController(
	usecase usecases.ChatUsecase,
	messageUsecase usecases.MessageUsecase,
) *chatController {
	return &chatController{
		usecase:        usecase,
		messageUsecase: messageUsecase,
	}
}

func (c *chatController) GETChats(ctx echo.Context) error {
	name, errorMessage := utils.GetStringQueryParam(ctx, "name")
	if errorMessage != "" {
		return ctx.JSON(http.StatusBadRequest, log.Response{
			Message: errorMessage,
		})
	}

	category, errorMessage := utils.GetStringQueryParam(ctx, "category")
	if errorMessage != "" {
		return ctx.JSON(http.StatusBadRequest, log.Response{
			Message: errorMessage,
		})
	}

	chats, err := c.usecase.List(dtos.ChatFilters{
		Name:     &name,
		Category: &category,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, log.Response{
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, chats)
}

func (uc *chatController) GETChatMessages(ctx echo.Context) error {
	id, errorMessage := utils.GetIntParam(ctx, "id")
	if errorMessage != "" {
		return ctx.JSON(http.StatusBadRequest, log.Response{
			Message: errorMessage,
		})
	}

	messages, err := uc.messageUsecase.ListByChatId(int64(id))
	if err != nil {
		if err == errors.ErrChatNotFound {
			return ctx.JSON(http.StatusBadRequest, log.Response{
				Message: err.Error(),
			})
		}

		return ctx.JSON(http.StatusInternalServerError, log.Response{
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, messages)
}

func (uc *chatController) POSTChat(ctx echo.Context) error {
	var body dtos.CreateChat

	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, log.Response{
			Message: err.Error(),
		})
	}

	err := validator.New().Struct(body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, log.Response{
			Message: err.Error(),
		})
	}

	chat, err := uc.usecase.Create(body)
	if err != nil {
		if err == errors.ErrChatNotFound {
			return ctx.JSON(http.StatusBadRequest, log.Response{
				Message: err.Error(),
			})
		}

		return ctx.JSON(http.StatusInternalServerError, log.Response{
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, chat)
}

func (uc *chatController) PUTChat(ctx echo.Context) error {
	id, errorMessage := utils.GetIntParam(ctx, "id")
	if errorMessage != "" {
		return ctx.JSON(http.StatusBadRequest, log.Response{
			Message: errorMessage,
		})
	}

	var body dtos.UpdateChat

	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, log.Response{
			Message: err.Error(),
		})
	}

	err := validator.New().Struct(body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, log.Response{
			Message: err.Error(),
		})
	}

	chat, err := uc.usecase.Update(int64(id), body)
	if err != nil {
		if err == errors.ErrChatNotFound {
			return ctx.JSON(http.StatusBadRequest, log.Response{
				Message: err.Error(),
			})
		}

		return ctx.JSON(http.StatusInternalServerError, log.Response{
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, chat)
}

func (uc *chatController) DELETEChat(ctx echo.Context) error {
	id, errorMessage := utils.GetIntParam(ctx, "id")
	if errorMessage != "" {
		return ctx.JSON(http.StatusBadRequest, log.Response{
			Message: errorMessage,
		})
	}

	err := uc.usecase.Delete(int64(id))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, log.Response{
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, log.Response{
		Message: "Chat deleted",
	})
}
