package handler

import (
	"fmt"
	"log"
	"net/http"
	"tservice-checker/internal/core"
	"tservice-checker/pkg"

	"github.com/labstack/echo/v4"
)

// * POST /api/sessions			— извлечение сессии из переданных данных от клиента (zip,tdata);
// @Summary ExtractSession
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
// @Router /api/sessions/extract [post]
func (h *Handler) ExtractSession() echo.HandlerFunc {
	return func(c echo.Context) error {
		pkg.InfoPrint("transport", "new request", "new request /api/sessions/extract")

		userID, err := getUserID(c)
		var totalResult = core.ExtractSessionResult{}

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

		for _, file := range files {
			fmt.Println("ok")
			// * сохранить полученный архив
			filePath, err := h.services.Session.SaveZip(file)
			if err != nil {
				fmt.Println("error save zip:: ", err)

				totalResult.SaveZipCounter--
				continue
			}
			totalResult.SaveZipCounter++

			// * разархивирировать из архива директорию tdata
			tdataPath, err := h.services.Session.Unzip(filePath)
			if err != nil {

				fmt.Println("error:: ", err)
				totalResult.UnZipCounter--
				continue
			}
			totalResult.UnZipCounter++
			fmt.Println("tdata path:: ", tdataPath)
			// * вытащить из директории tdata сессию
			sessions, err := h.services.Session.ExtractSession(tdataPath)
			if err != nil {
				fmt.Println("error session:: ", err)
			}
			for _, s := range sessions {
				// * проверить жива ли сессия
				// * если сессия жива сохранить ее базе
				h.services.Session.ValidateSession(&s)
				if err != nil {
					fmt.Println("error session:: ", err)
					log.Fatal()
				}
				// fmt.Println("session:: ", string(s.Data))
			}

		}

		return c.HTML(http.StatusOK, fmt.Sprintf("<p>Uploaded successfully %d files </p>", len(files)))

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
