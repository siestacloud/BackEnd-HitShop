package handler

import (
	"hitshop/pkg"
	"net/http"

	"github.com/labstack/echo/v4"
)

// * POST /api/save			— извлечение сессии из переданных данных от клиента (zip,tdata);
// @Summary MultipartSave
// @Security ApiKeyAuth
// @Tags Session
// @Description extract session by tdata folder, validate it and save in DB
// @ID extract-session
// @Accept  text/plain
// @Produce  text/plain
// @Param input body integer true "new title and description for item"
// @Success 200,202 {int} int "no content"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/sessions/save [post]
func (h *Handler) MultipartSave() echo.HandlerFunc {
	return func(c echo.Context) error {
		pkg.InfoPrint("transport", "new request", "new request /api/sessions/extract")

		userID, err := getUserID(c)

		if err != nil {
			pkg.ErrPrint("transport", http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, err.Error()) // в контексте нет id пользователя
		}

		pkg.InfoPrint("transport", "progress", userID)

		//* Multipart form
		form, err := c.MultipartForm()
		if err != nil {
			return err
		}
		files := form.File["files"]
		result, err := h.services.TAccount.MultipartSave(files)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error()) // ошибка при извлечении сессий или при сохр в бд
		}

		return c.JSON(http.StatusOK, result)
	}
}

// * POST /api/sessions/:phone  	— создание сессии по переданному номеру телефона (требует передачу проверочного кода);
// @Summary CreateSession
// @Security ApiKeyAuth
// @Tags Session
// @Description create session by phone number
// @ID create-session
// @Accept  json
// @Produce  json
// @Param phone path integer true "Phone number"
// @Success 200,202 {int} int "no content"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/sessions/{phone} [post]
func (h *Handler) CreateSession() echo.HandlerFunc {
	// todo реализовать
	return func(c echo.Context) error {
		pkg.InfoPrint("transport", "new request", "new request ", c.Request().RequestURI, " POST")
		return c.NoContent(http.StatusOK)

	}
}

// * GET /api/sessions/:phone  		— проверка наличия(и ее живучесть) сохраненной сессии по переданному номеру телефона;
// @Summary GetSession
// @Security ApiKeyAuth
// @Tags Session
// @Description get session by phone number
// @ID get-session
// @Accept  json
// @Produce  json
// @Param id path integer true "phone number"
// @Success 200 {object}  core.Session
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/sessions/{phone} [get]
func (h *Handler) GetSessionByPhoneNumber() echo.HandlerFunc {
	return func(c echo.Context) error {
		// todo реализовать ендпоинт
		pkg.InfoPrint("transport", "new request", "new request ", c.Request().RequestURI, " GET")
		return c.NoContent(http.StatusOK)

	}
}

// * DELETE /api/sessions/:phone  		— удаление сохраненной сессии по переданному номеру телефона;
// @Summary DeleteSession
// @Security ApiKeyAuth
// @Tags Session
// @Description delete session by phone number
// @ID delete_session-by-phone
// @Accept  json
// @Produce  json
// @Param phone path integer true "Phone number"
// @Success 200 {string} string "OK"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/sessions/{id} [delete]
func (h *Handler) DeleteSession() echo.HandlerFunc {
	return func(c echo.Context) error {
		// todo реализовать ендпоинт
		return c.JSON(http.StatusOK, statusResponse{Status: "ok"})

	}
}

// * PUT /api/sessions/:phone  		— обноление сохраненной сессии по переданному номеру телефона;
// @Summary UpdateSession
// @Security ApiKeyAuth
// @Tags Session
// @Description update session by ID
// @ID update_session_by_id
// @Accept  json
// @Produce  json
// @Param id path integer true "session ID in data base"
// @Param input body core.UpdateSessionInput true "new title and description"
// @Success 200 {string} string "OK"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/sessions/{id} [put]
func (h *Handler) UpdateSession() echo.HandlerFunc {
	return func(c echo.Context) error {
		// todo реализовать ендпоинт
		return c.JSON(http.StatusOK, statusResponse{"ok"})

	}
}
