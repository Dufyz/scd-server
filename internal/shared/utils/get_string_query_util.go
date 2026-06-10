package utils

import "github.com/labstack/echo/v4"

func GetStringQueryParam(ctx echo.Context, paramName string) (string, string) {
	param := ctx.QueryParam(paramName)

	if param == "" {
		return "", paramName + " is required!"
	}

	return param, ""
}
