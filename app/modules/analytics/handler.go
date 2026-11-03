package analytics

import (
	"fmt"
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

// PublicGetAnalytics provides unauthenticated analytics inspection for the web UI
func (h *Handler) PublicGetAnalytics(c echo.Context) error {
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

// ExportCSV exports raw click telemetry as a downloadable CSV file
func (h *Handler) ExportCSV(c echo.Context) error {
	code := strings.TrimSpace(c.Param("code"))
	if code == "" {
		return utils.JSONError(c, appErrors.BadRequest("missing short code parameter", nil))
	}

	csvBytes, err := h.service.ExportClicksCSV(c.Request().Context(), code)
	if err != nil {
		return utils.JSONError(c, err)
	}

	fileName := fmt.Sprintf("clicks_%s.csv", code)
	c.Response().Header().Set("Content-Type", "text/csv; charset=utf-8")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	return c.Blob(http.StatusOK, "text/csv; charset=utf-8", csvBytes)
}
