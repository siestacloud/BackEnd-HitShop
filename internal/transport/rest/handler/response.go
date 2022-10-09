package handler

import (
	"tservice-checker/pkg"

	"github.com/labstack/echo/v4"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func errResponse(c echo.Context, statusCode int, message string) error {
	pkg.ErrPrint("transport", statusCode, message)
	return c.JSON(statusCode, errorResponse{"error"})
}
