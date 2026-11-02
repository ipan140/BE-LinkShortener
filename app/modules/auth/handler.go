package auth

import (
	"net/http"

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

func (h *Handler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.JSONError(c, appErrors.BadRequest("invalid request payload", err.Error()))
	}

	resp, err := h.service.Register(c.Request().Context(), req)
	if err != nil {
		return utils.JSONError(c, err)
	}

	return utils.JSONSuccess(c, http.StatusCreated, "user registered successfully", resp)
}

func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.JSONError(c, appErrors.BadRequest("invalid request payload", err.Error()))
	}

	resp, err := h.service.Login(c.Request().Context(), req)
	if err != nil {
		return utils.JSONError(c, err)
	}

	return utils.JSONSuccess(c, http.StatusOK, "login successful", resp)
}

func (h *Handler) Me(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return utils.JSONError(c, appErrors.Unauthorized("unauthorized"))
	}

	user, err := h.service.GetProfile(c.Request().Context(), userID)
	if err != nil {
		return utils.JSONError(c, err)
	}

	return utils.JSONSuccess(c, http.StatusOK, "user profile retrieved", user)
}
