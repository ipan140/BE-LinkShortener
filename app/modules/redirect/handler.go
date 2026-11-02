package redirect

import (
	"fmt"
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

// Redirect handles ultra-fast redirection, UTM forwarding, and password challenges
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

	// Password can be supplied via query (?p=... or ?password=...) or form
	password := c.QueryParam("p")
	if password == "" {
		password = c.QueryParam("password")
	}
	if password == "" && req.Method == http.MethodPost {
		password = c.FormValue("password")
	}

	incomingQuery := req.URL.RawQuery

	click := &ClickData{
		ShortCode:   code,
		IP:          utils.GetRealIP(c),
		Country:     country,
		City:        city,
		UserAgent:   uaString,
		DeviceType:  deviceInfo.DeviceType,
		Browser:     deviceInfo.Browser,
		OS:          deviceInfo.OS,
		Referrer:    req.Referer(),
		QueryParams: incomingQuery,
		ClickedAt:   time.Now(),
	}

	targetURL, requiresPassword, err := h.service.Resolve(req.Context(), code, incomingQuery, password, click)
	if err != nil {
		return utils.JSONError(c, err)
	}

	// If the link is protected and no valid password was provided yet
	if requiresPassword {
		accept := req.Header.Get("Accept")
		if strings.Contains(accept, "application/json") {
			return utils.JSONError(c, appErrors.New(http.StatusUnauthorized, "PASSWORD_REQUIRED", "this link is password protected", map[string]string{
				"short_code": code,
			}))
		}

		// Render responsive password challenge page
		html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Protected Link</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; background: #0f172a; color: #f8fafc; }
        .card { background: #1e293b; padding: 2rem; border-radius: 12px; box-shadow: 0 10px 25px rgba(0,0,0,0.5); width: 100%%; max-width: 360px; text-align: center; }
        h2 { margin-top: 0; font-size: 1.25rem; font-weight: 600; color: #38bdf8; }
        p { color: #94a3b8; font-size: 0.875rem; margin-bottom: 1.5rem; }
        input[type="password"] { width: 100%%; padding: 0.75rem; border: 1px solid #334155; border-radius: 6px; background: #0f172a; color: #fff; box-sizing: border-box; margin-bottom: 1rem; font-size: 1rem; outline: none; }
        input[type="password"]:focus { border-color: #38bdf8; }
        button { width: 100%%; padding: 0.75rem; background: #0284c7; color: white; border: none; border-radius: 6px; font-weight: 600; font-size: 1rem; cursor: pointer; transition: background 0.2s; }
        button:hover { background: #0369a1; }
    </style>
</head>
<body>
    <div class="card">
        <h2>🔒 Protected Link</h2>
        <p>This destination is password protected. Enter the access key to proceed.</p>
        <form method="POST" action="/%s">
            <input type="password" name="password" placeholder="Enter password" required autofocus />
            <button type="submit">Unlock & Continue</button>
        </form>
    </div>
</body>
</html>`, code)
		return c.HTML(http.StatusOK, html)
	}

	// HTTP 302 Found redirect with preserved UTMs
	return c.Redirect(http.StatusFound, targetURL)
}
