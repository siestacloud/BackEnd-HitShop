package handler

import (
	"hitshop/pkg"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// * GET /auth/static/register			— отдает статику с form-data для регистрации пользователя
// @Summary StaticRegister
// @Tags Static
// @Description ендпоинт отдает статику с form-data для регистрации пользователя
// @ID static-register
// @Accept  text/plain
// @Produce  text/plain
// @Success 200,202 {string} string "html"
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/static/register [get]
func (h *Handler) StaticRegister() echo.HandlerFunc {
	return func(c echo.Context) error {
		pkg.InfoPrint("transport", "process", "new request /auth/static/register")
		html, err := os.ReadFile("./assets/register.html")
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}
		return c.HTML(http.StatusOK, string(html))
	}
}

// * GET /auth/static/login		— отдает статику с form-data для аутентификации пользователя
// @Summary StaticLogin
// @Tags Static
// @Description ендпоинт отдает статику с form-data для аутентификации пользователя
// @ID static-login
// @Accept  text/plain
// @Produce  text/plain
// @Success 200,202 {string} string "html"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/static/login [get]
func (h *Handler) StaticLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		pkg.InfoPrint("transport", "process", "new request /auth/static/login")
		html, err := os.ReadFile("./assets/login.html")
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}
		return c.HTML(http.StatusOK, string(html))
	}
}

// * GET /			— отдает статику с MultipartForm-data для передачи архивов от клиента
// @Summary StaticExtract
// @Security ApiKeyAuth
// @Tags Static
// @Description ендпоинт отдает статику с form-data для передачи архивов от клиента
// @ID static-extract
// @Accept  text/plain
// @Produce  text/plain
// @Success 200,202 {string} string "html"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/sessions/static/extract [get]
func (h *Handler) StaticExtract() echo.HandlerFunc {
	return func(c echo.Context) error {
		pkg.InfoPrint("transport", "process", "new request /api/sessions/static/extract GET")
		userID, err := getUserID(c)
		if err != nil {
			pkg.ErrPrint("transport", http.StatusInternalServerError, err)
			return c.Redirect(http.StatusSeeOther, "/auth/static/login")
		}
		pkg.InfoPrint("transport", "progress", "user: ", userID, " access allowed")

		html, err := os.ReadFile("./assets/extract.html")
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}
		return c.HTML(http.StatusOK, string(html))
	}
}

// * GET /			— перенаправляет все запросы на главную страницу
// @Summary RedirectToExtract
// @Tags Static
// @Description ендпоинт перенаправляет все запросы на основную страницу
// @ID static-redirect
// @Accept  text/plain
// @Produce  text/plain
// @Success 200,202 {int} int "no content"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router / [get]
func (h *Handler) RedirectToExtract() echo.HandlerFunc {
	return func(c echo.Context) error {
		pkg.InfoPrint("transport", "process", "new request on "+c.Request().RequestURI+" will redirect to /api/sessions/static/extract")
		return c.Redirect(http.StatusSeeOther, "/api/sessions/static/extract")
	}
}
