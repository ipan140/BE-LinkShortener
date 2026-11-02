package redirect

import (
	"net/http"
	"strings"
	"time"

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

// Redirect handles the high-performance redirect endpoint
func (h *Handler) Redirect(c echo.Context) error {
	code := strings.TrimSpace(c.Param("code"))
	if code == "" {
		return utils.JSONError(c, appErrors.BadRequest("missing short code", nil))
	}

	req := c.Request()
	uaString := req.UserAgent()
	deviceInfo := utils.ParseUserAgent(uaString)

	country := req.Header.Get("CF-IPCountry")
	if country == "" {
		country = req.Header.Get("X-Country-Code")
	}
	if country == "" {
		country = "Unknown"
	}

	city := req.Header.Get("CF-IPCity")
	if city == "" {
		city = req.Header.Get("X-City")
	}
	if city == "" {
		city = "Unknown"
	}

	click := &ClickData{
		ShortCode:  code,
		IP:         utils.GetRealIP(c),
		Country:    country,
		City:       city,
		UserAgent:  uaString,
		DeviceType: deviceInfo.DeviceType,
		Browser:    deviceInfo.Browser,
		OS:         deviceInfo.OS,
		Referrer:   req.Referer(),
		ClickedAt:  time.Now(),
	}

	targetURL, err := h.service.Resolve(req.Context(), code, click)
	if err != nil {
		return utils.JSONError(c, err)
	}

	// HTTP 302 Found (or 307 Temporary Redirect) to ensure subsequent clicks are tracked
	return c.Redirect(http.StatusFound, targetURL)
}
