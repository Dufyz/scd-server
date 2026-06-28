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

type messageController struct {
	usecase  usecases.MessageUsecase
	validate *validator.Validate
}

func NewMessageController(
	usecase usecases.MessageUsecase,
) *messageController {
	return &messageController{
		usecase:  usecase,
		validate: validator.New(),
	}
}

func (uc *messageController) POSTMessage(ctx echo.Context) error {
	var body dtos.CreateMessage

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

	message, err := uc.usecase.Create(body)
	if err != nil {
		if err == errors.ErrMessageFKChatId {
			return ctx.JSON(http.StatusBadRequest, log.Response{
				Message: err.Error(),
			})
		}

		return ctx.JSON(http.StatusInternalServerError, log.Response{
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, message)
}

func (uc *messageController) PUTMessage(ctx echo.Context) error {
	id, errorMessage := utils.GetIntParam(ctx, "id")
	if errorMessage != "" {
		return ctx.JSON(http.StatusBadRequest, log.Response{
			Message: errorMessage,
		})
	}

	var body dtos.UpdateMessage

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

	message, err := uc.usecase.Update(int64(id), body)
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

	return ctx.JSON(http.StatusOK, message)
}

func (uc *messageController) DELETEMessage(ctx echo.Context) error {
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
		Message: "Message deleted",
	})
}
