package handler

import (
	"hitshop/internal/core"
	"hitshop/pkg"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statResponse struct {
	Status string `json:"status"`
}

func errResponse(c echo.Context, statusCode int, message string) error {
	pkg.ErrPrintT(c.Request().RequestURI, statusCode, message)
	return c.JSON(statusCode, errorResponse{"error"})
}

func Redirect(c echo.Context, statusCode int, message string, URI string) error {
	pkg.WarnPrintT(c.Request().RequestURI, statusCode, message)
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

func accountResponse(c echo.Context, statusCode int, acc core.Account) error {
	pkg.InfoPrintT(c.Request().RequestURI, statusCode, acc.Email)
	return c.JSON(statusCode, acc)
}

func statusResponse(c echo.Context, statusCode int, message string) error {
	pkg.InfoPrintT(c.Request().RequestURI, statusCode, message)
	return c.JSON(statusCode, statResponse{message})
}
