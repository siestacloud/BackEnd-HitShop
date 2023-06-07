package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	authorizationCookie = "Token"
	userCtx             = "userID"
)

func (h *Handler) UserIdentity(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		cookie, err := c.Cookie(authorizationCookie)
		if err != nil {
			return Redirect(c, http.StatusSeeOther, err.Error()+" client will redirect", "/auth/static/login")
		}

		headerParts := strings.Split(cookie.Value, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return Redirect(c, http.StatusSeeOther, "invalid cookie, client will redirect", "/auth/static/login")
		}

		if len(headerParts[1]) == 0 {
			return Redirect(c, http.StatusSeeOther, "token in cookie is empty, client will redirect", "/auth/static/login")
		}

		userID, err := h.services.Authorization.ParseToken(headerParts[1])
		if err != nil {
			return Redirect(c, http.StatusSeeOther, err.Error()+" client will redirect", "/auth/static/login")
		}
		// Добавляю ID пользователя в контекст
		c.Set(userCtx, userID)
		return next(c)
	}
}

func getUserID(c echo.Context) (uuid.UUID, error) {
	id := c.Get(userCtx)

	idInt, ok := id.(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("user id is of invalid type")
	}
	return idInt, nil
}

func (h *Handler) CheckContentType(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Content-Type") != "application/json" {
			return errResponse(c, http.StatusBadRequest, "refuse request")
		}
		return next(c)
	}
}
