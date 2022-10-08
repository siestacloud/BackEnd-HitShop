package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// * POST /api/sessions			— извлечение сессии из переданных данных от клиента (zip,tdata);
// @Summary SignUp
// @Tags session
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body core.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/user/register [post]
func (h *Handler) ExtractSession() echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

// * POST /api/sessions/:phone  	— создание сессии по переданному номеру телефона (требует передачу проверочного кода);
// @Summary SignUp
// @Tags session
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body core.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/user/register [post]
func (h *Handler) CreateSession() echo.HandlerFunc {

	return func(c echo.Context) error {
		return c.NoContent(http.StatusOK)

	}
}

// * GET /api/sessions/:phone  		— проверка наличия(и ее живучесть) сохраненной сессии по переданному номеру телефона;
// @Summary GetSession
// @Security ApiKeyAuth
// @Tags session
// @Description get list by id
// @ID get_list_by_id
// @Accept  json
// @Produce  json
// @Param id path integer true "phone number"
// @Success 200 {object}  core.Session
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/{id} [get]
func (h *Handler) GetSessionByPhoneNumber() echo.HandlerFunc {

	return func(c echo.Context) error {
		return c.NoContent(http.StatusOK)

	}
}
