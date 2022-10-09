package handler

import (
	"net/http"
	"time"
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

func Redirect(c echo.Context, statusCode int, message string, URI string) error {
	pkg.WarnPrint("transport", statusCode, message)
	return c.Redirect(statusCode, URI)
}

// writeCookie Создаю новую куку
func writeCookie(path, name, value string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Path = path
	cookie.Expires = time.Now().Add(24 * time.Hour)
	return cookie
}
