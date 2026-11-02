package analytics

import (
	"net/http"
	"strings"

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

func (h *Handler) GetAnalytics(c echo.Context) error {
	code := strings.TrimSpace(c.Param("code"))
	if code == "" {
		return utils.JSONError(c, appErrors.BadRequest("missing short code parameter", nil))
	}

	summary, err := h.service.GetAnalytics(c.Request().Context(), code)
	if err != nil {
		return utils.JSONError(c, err)
	}

	return utils.JSONSuccess(c, http.StatusOK, "analytics retrieved successfully", summary)
}
