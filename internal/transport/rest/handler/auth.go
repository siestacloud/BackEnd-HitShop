package handler

import (
	"fmt"
	"net/http"
	"strings"
	"tservice-checker/internal/core"
	"tservice-checker/pkg"

	"github.com/labstack/echo/v4"
)

// 	* `POST /auth/register` 	— регистрация пользователя;
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
		pkg.InfoPrint("transport", "process", "new request on /auth/register")

		var input core.User
		// todo обработку fetch запроса с данными в формате json
		//? if err := c.Bind(&input); err != nil {
		//? 	pkg.ErrPrint("transport", http.StatusBadRequest, err)
		//? 	return errResponse(c, http.StatusBadRequest, "bind body failure")
		//? }
		input.Login = c.FormValue("login")
		input.Password = c.FormValue("password")
		fmt.Println("INPUT ", input)
		if err := c.Validate(input); err != nil {
			pkg.ErrPrint("transport", http.StatusBadRequest, err)
			return errResponse(c, http.StatusBadRequest, "validate failure")
		}

		// * авторизация
		_, err := h.services.Authorization.CreateUser(input)
		if err != nil {
			if strings.Contains(err.Error(), "login busy") {
				pkg.ErrPrint("transport", http.StatusConflict, err)
				return errResponse(c, http.StatusConflict, err.Error())
			}

			pkg.ErrPrint("transport", http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, "internal server error")
		}

		// * аутентификация
		token, err := h.services.Authorization.GenerateToken(input.Login, input.Password)
		if err != nil {

			pkg.ErrPrint("transport", http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, "internal server error")
		}

		c.SetCookie(writeCookie("/", "Token", "Bearer "+token))
		return c.Redirect(http.StatusSeeOther, "/api/sessions/static/extract")
	}
}

type signInInput struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// 	* `POST /auth/login` 						— аутентификация пользователя;
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
		pkg.InfoPrint("transport", "process", "new request on /auth/login")
		var input signInInput
		// todo обработку fetch запроса с данными в формате json
		//? if err := c.Bind(&input); err != nil {
		//? 	pkg.ErrPrint("transport", http.StatusBadRequest, err)
		//? 	return errResponse(c, http.StatusBadRequest, "bind body failure")
		//? }
		input.Login = c.FormValue("login")
		input.Password = c.FormValue("password")
		if err := c.Validate(input); err != nil {
			pkg.ErrPrint("transport", http.StatusBadRequest, err)
			return errResponse(c, http.StatusBadRequest, "validate failure")
		}
		token, err := h.services.Authorization.GenerateToken(input.Login, input.Password)
		if err != nil {
			if strings.Contains(err.Error(), "invalid username/password pair") {
				pkg.ErrPrint("transport", http.StatusBadRequest, err)
				return errResponse(c, http.StatusUnauthorized, err.Error())
			}

			pkg.ErrPrint("transport", http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, "internal server error")
		}

		c.SetCookie(writeCookie("/", "Token", "Bearer "+token))
		return c.Redirect(http.StatusMovedPermanently, "/api/sessions/static/extract")
	}
}
