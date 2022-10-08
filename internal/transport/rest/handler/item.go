package handler

import (
	"gitlab.com/siteasservice/project-architecture/templates/template-svc-golang/internal/core"

	"github.com/labstack/echo/v4"

	"net/http"
	"strconv"
)

// @Summary CreateItem
// @Security ApiKeyAuth
// @Tags Item
// @Description create item in list
// @ID create_item_in_list
// @Accept  json
// @Produce  json
// @Param id path integer true "List ID"
// @Param input body core.TodoItem true "new title and description for item"
// @Success 200 {int} int "item id"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/{id}/items/ [post]
func (h *Handler) CreateItem() echo.HandlerFunc {

	return func(c echo.Context) error {
		userId, err := getUserId(c)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}

		listId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errResponse(c, http.StatusBadRequest, "invalid list id param")
		}

		var input core.TodoItem
		if err := c.Bind(&input); err != nil {
			return errResponse(c, http.StatusBadRequest, err.Error())
		}

		id, err := h.services.TodoItem.Create(userId, listId, input)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"id": id,
		})
	}
}

// @Summary GetAllItems
// @Security ApiKeyAuth
// @Tags Item
// @Description get all items
// @ID get_all_items
// @Accept  json
// @Produce  json
// @Param id path integer true "List ID"
// @Success 200 {object} []core.TodoItem "Data items list"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/{id}/items/ [get]
func (h *Handler) GetAllItems() echo.HandlerFunc {

	return func(c echo.Context) error {
		userId, err := getUserId(c)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}

		listId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errResponse(c, http.StatusBadRequest, "invalid list id param")
		}

		items, err := h.services.TodoItem.GetAll(userId, listId)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, items)
	}
}

// @Summary GetItemById
// @Security ApiKeyAuth
// @Tags Item
// @Description get item by id
// @ID get_item_by_id
// @Accept  json
// @Produce  json
// @Param id path integer true "Item ID"
// @Success 200 {object} core.TodoItem "Data item"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/items/{id} [get]
func (h *Handler) GetItemById() echo.HandlerFunc {

	return func(c echo.Context) error {
		userId, err := getUserId(c)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}

		itemId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errResponse(c, http.StatusBadRequest, "invalid list id param")
		}

		item, err := h.services.TodoItem.GetById(userId, itemId)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, item)
	}
}

// @Summary UpdateItem
// @Security ApiKeyAuth
// @Tags Item
// @Description update item by id
// @ID update_item_by_id
// @Accept  json
// @Produce  json
// @Param id path integer true "Item ID"
// @Param input body core.UpdateItemInput true "new title and description for item"
// @Success 200 {string} string "OK"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/items/{id} [put]
func (h *Handler) UpdateItem() echo.HandlerFunc {

	return func(c echo.Context) error {
		userId, err := getUserId(c)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errResponse(c, http.StatusBadRequest, "invalid id param")
		}

		var input core.UpdateItemInput
		if err := c.Bind(&input); err != nil {
			return errResponse(c, http.StatusBadRequest, err.Error())
		}

		if err := h.services.TodoItem.Update(userId, id, input); err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, statusResponse{"ok"})
	}
}

// @Summary DeleteItem
// @Security ApiKeyAuth
// @Tags Item
// @Description delete item by id
// @ID delete_item_by_id
// @Accept  json
// @Produce  json
// @Param id path integer true "Item ID"
// @Success 200 {string} string "OK"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/items/{id} [delete]
func (h *Handler) DeleteItem() echo.HandlerFunc {

	return func(c echo.Context) error {
		userId, err := getUserId(c)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}

		itemId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errResponse(c, http.StatusBadRequest, "invalid list id param")
		}

		err = h.services.TodoItem.Delete(userId, itemId)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, statusResponse{"ok"})
	}
}
