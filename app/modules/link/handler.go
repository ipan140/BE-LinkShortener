package link

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	appErrors "link-shortener/common/errors"
	"link-shortener/common/utils"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return utils.JSONError(c, appErrors.Unauthorized("unauthorized"))
	}

	var req CreateLinkRequest
	if err := c.Bind(&req); err != nil {
		return utils.JSONError(c, appErrors.BadRequest("invalid request payload", err.Error()))
	}

	resp, err := h.service.CreateLink(c.Request().Context(), userID, req)
	if err != nil {
		return utils.JSONError(c, err)
	}

	return utils.JSONSuccess(c, http.StatusCreated, "short link created successfully", resp)
}

func (h *Handler) GetByID(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return utils.JSONError(c, appErrors.Unauthorized("unauthorized"))
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return utils.JSONError(c, appErrors.BadRequest("invalid id parameter", nil))
	}

	resp, err := h.service.GetLinkByID(c.Request().Context(), uint(id), userID)
	if err != nil {
		return utils.JSONError(c, err)
	}

	return utils.JSONSuccess(c, http.StatusOK, "link retrieved successfully", resp)
}

func (h *Handler) List(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return utils.JSONError(c, appErrors.Unauthorized("unauthorized"))
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	links, total, err := h.service.ListLinks(c.Request().Context(), userID, page, limit)
	if err != nil {
		return utils.JSONError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"data":    links,
		"meta": map[string]any{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func (h *Handler) Update(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return utils.JSONError(c, appErrors.Unauthorized("unauthorized"))
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return utils.JSONError(c, appErrors.BadRequest("invalid id parameter", nil))
	}

	var req UpdateLinkRequest
	if err := c.Bind(&req); err != nil {
		return utils.JSONError(c, appErrors.BadRequest("invalid request payload", err.Error()))
	}

	resp, err := h.service.UpdateLink(c.Request().Context(), uint(id), userID, req)
	if err != nil {
		return utils.JSONError(c, err)
	}

	return utils.JSONSuccess(c, http.StatusOK, "link updated successfully", resp)
}

func (h *Handler) Delete(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return utils.JSONError(c, appErrors.Unauthorized("unauthorized"))
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return utils.JSONError(c, appErrors.BadRequest("invalid id parameter", nil))
	}

	if err := h.service.DeleteLink(c.Request().Context(), uint(id), userID); err != nil {
		return utils.JSONError(c, err)
	}

	return utils.JSONSuccess(c, http.StatusOK, "link deleted successfully", nil)
}
