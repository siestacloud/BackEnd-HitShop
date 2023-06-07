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

		// pkg.InfoPrintT(uri, "healfy", "detect request")
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

		// Generate Verification Code
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

		// c.SetCookie(writeCookie("/", "Token", "Bearer "+token))
		return statusResponse(c, http.StatusOK, "success")

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

		pkg.InfoPrintT(uri, "process", "detect request")
		// var acc core.Account
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
		return c.JSON(http.StatusOK, "login success")
	}
}

// Logout
func (h *Handler) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.SetCookie(writeCookie("/", "Token", ""))
		return c.JSON(http.StatusOK, "logout success")
	}
}

// VerifyEmail
func (h *Handler) VerifyEmail() echo.HandlerFunc {
	return func(c echo.Context) error {
		// uri := c.Request().RequestURI
		code := c.Param("code")
		verification_code := pkg.Encode(code)

		acc, err := h.services.GetAccountByCode(verification_code)
		if err != nil {
			return errResponse(c, http.StatusConflict, err.Error())
		}
		if acc.Verified {
			return errResponse(c, http.StatusConflict, fmt.Sprintf("Email %s already verified", acc.Email))
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
