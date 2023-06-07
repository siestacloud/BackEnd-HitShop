package handler

import (
	"fmt"
	"hitshop/internal/core"
	"hitshop/pkg"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/thanhpk/randstr"
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
		var payload core.SignUpInput
		if err := c.Bind(&payload); err != nil {
			pkg.ErrPrintT(uri, "error", err)
			return errResponse(c, http.StatusBadRequest, "body bind failure")
		}
		if err := c.Validate(payload); err != nil {
			pkg.ErrPrintT(uri, http.StatusBadRequest, err)
			return errResponse(c, http.StatusBadRequest, "body validate failure")
		}
		if payload.Password != payload.PasswordConfirm {
			return errResponse(c, http.StatusBadRequest, "passwords do not match")
		}

		//* Generate Verification Code
		code := randstr.String(20)
		verification_code := pkg.Encode(code)

		now := time.Now()
		acc := core.Account{
			Email:            strings.ToLower(payload.Email),
			Password:         payload.Password,
			Role:             "client",
			Verified:         false,
			Status:           "healfy",
			VerificationCode: verification_code,
			CreateAt:         now,
			UpdateAt:         now,
		}

		// * авторизация
		_, err := h.services.Authorization.CreateAccount(&acc)
		if err != nil {
			if strings.Contains(err.Error(), "login busy") {
				return errResponse(c, http.StatusConflict, err.Error())
			}
			pkg.ErrPrintT(uri, http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, "internal server error")
		}
		if err := pkg.SendEmail(&acc, code, h.cfg); err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}
		return statusResponse(c, http.StatusOK, fmt.Sprintf("Success send email to %s", acc.Email))

	}
}

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

		var payload core.SignInInput
		if err := c.Bind(&payload); err != nil {
			pkg.ErrPrintT(uri, http.StatusBadRequest, err)
			return errResponse(c, http.StatusBadRequest, "bind body failure")
		}
		if err := c.Validate(payload); err != nil {
			pkg.ErrPrintT(uri, http.StatusBadRequest, err)
			return errResponse(c, http.StatusBadRequest, "validate failure")
		}
		token, err := h.services.Authorization.GenerateToken(payload.Email, payload.Password)
		if err != nil {
			if strings.Contains(err.Error(), "invalid username/password pair") {
				return errResponse(c, http.StatusUnauthorized, err.Error())
			}
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}
		c.SetCookie(writeCookie("/", "Token", "Bearer "+token))
		return statusResponse(c, http.StatusOK, fmt.Sprintf("Success login client %s", payload.Email))
	}
}

// Logout удаление cookie
func (h *Handler) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.SetCookie(writeCookie("/", "Token", ""))
		return statusResponse(c, http.StatusOK, "Success logout client")
	}
}

// VerifyEmail подтверждение почты клиента
func (h *Handler) VerifyEmail() echo.HandlerFunc {
	return func(c echo.Context) error {
		code := c.Param("code")
		verification_code := pkg.Encode(code)

		acc, err := h.services.GetAccountByCode(verification_code)
		if err != nil {
			return errResponse(c, http.StatusConflict, err.Error())
		}
		if acc.Verified {
			return Redirect(c, http.StatusSeeOther, fmt.Sprintf("Email %s already verified, Redirect to / ", acc.Email), "/")
		}
		acc.Verified = true
		acc.UpdateAt = time.Now()
		_, err = h.services.UpdateAccount(acc)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}

		return Redirect(c, http.StatusSeeOther, fmt.Sprintf("Email %s verified successfully, Redirect to /login ", acc.Email), "/login")
	}
}

// ChangePassword изменить пароль клиента
func (h *Handler) ChangePassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload core.ChangePassInput
		if err := c.Bind(&payload); err != nil {
			return errResponse(c, http.StatusBadRequest, err.Error())
		}
		if err := c.Validate(payload); err != nil {
			return errResponse(c, http.StatusBadRequest, err.Error())
		}
		accountUUID, err := getUserID(c)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error()) // в контексте нет id пользователя
		}
		_, err = h.services.ChangePassword(accountUUID, payload.Password, payload.PasswordNew)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}
		return statusResponse(c, http.StatusOK, "Success change password")

	}
}
