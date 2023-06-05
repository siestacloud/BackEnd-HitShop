package handler

import (
	"hitshop/internal/core"
	"hitshop/pkg"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

//   - `POST /auth/register` 	— регистрация пользователя;
//
// @Summary Register
// @Tags Auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body core.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/register [post]
func (h *Handler) Register() echo.HandlerFunc {

	return func(c echo.Context) error {
		uri := c.Request().RequestURI

		pkg.InfoPrint(uri, "ok", "detect request")
		var acc core.Account
		// todo обработку fetch запроса с данными в формате json
		if err := c.Bind(&acc); err != nil {
			pkg.ErrPrintT(uri, "error", err)
			return errResponse(c, http.StatusBadRequest, "body bind failure")
		}
		if err := c.Validate(acc); err != nil {
			pkg.ErrPrintT(uri, http.StatusBadRequest, err)
			return errResponse(c, http.StatusBadRequest, "body validate failure")
		}

		// * авторизация
		_, err := h.services.Authorization.CreateUser(acc)
		if err != nil {
			if strings.Contains(err.Error(), "login busy") {
				return errResponse(c, http.StatusConflict, err.Error())
			}

			pkg.ErrPrintT(uri, http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, "internal server error")
		}

		// * аутентификация
		token, err := h.services.Authorization.GenerateToken(acc.Email, acc.Password)
		if err != nil {
			pkg.ErrPrintT(uri, http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, "internal server error")
		}

		c.SetCookie(writeCookie("/", "Token", "Bearer "+token))
		return c.JSON(http.StatusOK, "register success")

	}
}

//// type signInInput struct {
//// 	Login    string `json:"login" validate:"required"`
//// 	Password string `json:"password" validate:"required"`
//// }

//   - `POST /auth/login` 						— аутентификация пользователя;
//
// @Summary Login
// @Tags Auth
// @Description login in account
// @ID login
// @Accept  json
// @Produce  json
// @Param input body signInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/login [post]
func (h *Handler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		uri := c.Request().RequestURI

		pkg.InfoPrint(uri, "process", "detect request")
		var acc core.Account
		if err := c.Bind(&acc); err != nil {
			pkg.ErrPrintT(uri, http.StatusBadRequest, err)
			return errResponse(c, http.StatusBadRequest, "bind body failure")
		}
		if err := c.Validate(acc); err != nil {
			pkg.ErrPrintT(uri, http.StatusBadRequest, err)
			return errResponse(c, http.StatusBadRequest, "validate failure")
		}
		token, err := h.services.Authorization.GenerateToken(acc.Email, acc.Password)
		if err != nil {
			if strings.Contains(err.Error(), "invalid username/password pair") {
				pkg.ErrPrintT(uri, http.StatusBadRequest, err)
				return errResponse(c, http.StatusUnauthorized, err.Error())
			}

			pkg.ErrPrintT(uri, http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, "internal server error")
		}

		c.SetCookie(writeCookie("/", "Token", "Bearer "+token))
		return c.JSON(http.StatusOK, "login success")

	}
}
