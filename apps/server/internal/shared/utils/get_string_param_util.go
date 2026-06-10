package utils

import "github.com/labstack/echo/v4"

func GetStringParam(ctx echo.Context, paramName string) (string, string) {
	param := ctx.Param(paramName)
	if param == "" {
		return "", paramName + " is required!"
	}

	return param, ""
}
