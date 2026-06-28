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
	validate       *validator.Validate
}

func NewChatController(
	usecase usecases.ChatUsecase,
	messageUsecase usecases.MessageUsecase,
) *chatController {
	return &chatController{
		usecase:        usecase,
		messageUsecase: messageUsecase,
		validate:       validator.New(),
	}
}

func (c *chatController) GETChats(ctx echo.Context) error {
	nameParam := ctx.QueryParam("name")
	categoryParam := ctx.QueryParam("category")

	var name *string
	if nameParam != "" {
		name = &nameParam
	}

	var category *string
	if categoryParam != "" {
		category = &categoryParam
	}

	chats, err := c.usecase.List(dtos.ChatFilters{
		Name:     name,
		Category: category,
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

	if err := uc.validate.Struct(body); err != nil {
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

	if err := uc.validate.Struct(body); err != nil {
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
