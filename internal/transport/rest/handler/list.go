package handler

import (
	"net/http"
	"strconv"

	"gitlab.com/siteasservice/project-architecture/templates/template-svc-golang/internal/core"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// @Summary CreateList
// @Security ApiKeyAuth
// @Tags List
// @Description create list
// @ID create_list
// @Accept  json
// @Produce  json
// @Param input body core.TodoList true "ToDo list"
// @Success 200 {string} string "list id"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/ [post]
func (h *Handler) CreateList() echo.HandlerFunc {

	return func(c echo.Context) error {
		userId, err := getUserId(c)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error()) // в контексте нет id пользователя
		}

		var input core.TodoList // Получаем данные от пользователя для создания списка, id создает соответствие список-пользователь
		if err := c.Bind(&input); err != nil {
			return errResponse(c, http.StatusBadRequest, err.Error())
		}
		logrus.Info(input)
		id, err := h.services.TodoList.Create(userId, input)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, map[string]interface{}{ //Возвращаю id списка
			"id": id,
		})
	}
}

type getAllListsResponse struct {
	Data []core.TodoList `json:"data"`
}

// @Summary GetAllLists
// @Security ApiKeyAuth
// @Tags List
// @Description get all  list
// @ID get_all_list
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Data list lists"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/ [get]
func (h *Handler) GetAllLists() echo.HandlerFunc {

	return func(c echo.Context) error {
		userId, err := getUserId(c)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}

		lists, err := h.services.TodoList.GetAll(userId)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, getAllListsResponse{
			Data: lists,
		})
	}
}

// @Summary GetListById
// @Security ApiKeyAuth
// @Tags List
// @Description get list by id
// @ID get_list_by_id
// @Accept  json
// @Produce  json
// @Param id path integer true "User ID"
// @Success 200 {object}  core.TodoList
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/{id} [get]
func (h *Handler) GetListById() echo.HandlerFunc {

	return func(c echo.Context) error {
		userId, err := getUserId(c)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errResponse(c, http.StatusBadRequest, "invalid id param")
		}

		list, err := h.services.TodoList.GetById(userId, id)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, list)
	}
}

// @Summary UpdateList
// @Security ApiKeyAuth
// @Tags List
// @Description update list by id
// @ID update_list_by_id
// @Accept  json
// @Produce  json
// @Param id path integer true "List ID"
// @Param input body core.UpdateListInput true "new title and description"
// @Success 200 {string} string "OK"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/{id} [put]
func (h *Handler) UpdateList() echo.HandlerFunc {

	return func(c echo.Context) error {
		userId, err := getUserId(c)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errResponse(c, http.StatusBadRequest, "invalid id param")
		}

		var input core.UpdateListInput
		if err := c.Bind(&input); err != nil {
			return errResponse(c, http.StatusBadRequest, err.Error())
		}

		if err := h.services.TodoList.Update(userId, id, input); err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, statusResponse{"ok"})

	}
}

// @Summary DeleteList
// @Security ApiKeyAuth
// @Tags List
// @Description delete list by id
// @ID delete_list_by_id
// @Accept  json
// @Produce  json
// @Param id path integer true "List ID"
// @Success 200 {string} string "OK"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/{id} [delete]
func (h *Handler) DeleteList() echo.HandlerFunc {

	return func(c echo.Context) error {
		userId, err := getUserId(c)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errResponse(c, http.StatusBadRequest, "invalid id param")
		}

		err = h.services.TodoList.Delete(userId, id)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, statusResponse{
			Status: "ok",
		})

	}
}
