package redirect

import (
	"fmt"
	"html"
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

func isSocialCrawler(ua string) bool {
	uaLower := strings.ToLower(ua)
	return strings.Contains(uaLower, "facebookexternalhit") ||
		strings.Contains(uaLower, "whatsapp") ||
		strings.Contains(uaLower, "telegrambot") ||
		strings.Contains(uaLower, "twitterbot") ||
		strings.Contains(uaLower, "discordbot") ||
		strings.Contains(uaLower, "slackbot") ||
		strings.Contains(uaLower, "linkedinbot") ||
		strings.Contains(uaLower, "pinterest")
}

// Redirect handles ultra-fast redirection, social previews, UTM forwarding, and device routing
func (h *Handler) Redirect(c echo.Context) error {
	code := strings.TrimSpace(c.Param("code"))
	if code == "" {
		return utils.JSONError(c, appErrors.BadRequest("missing short code", nil))
	}

	req := c.Request()
	uaString := req.UserAgent()

	// 1. Social Media Bot Detection (WhatsApp, Telegram, Twitter, Discord, etc.)
	if isSocialCrawler(uaString) {
		link, err := h.service.GetLinkDetails(req.Context(), code)
		if err == nil && link != nil && (link.Title != "" || link.Description != "" || link.ImageURL != "") {
			return h.renderSocialPreview(c, link)
		}
	}

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
		accept := req.Header.Get("Accept")
		if strings.Contains(accept, "text/html") && !strings.Contains(accept, "application/json") {
			return h.renderErrorPage(c, code, err)
		}
		return utils.JSONError(c, err)
	}

	// If the link is protected and requires a password
	if requiresPassword {
		accept := req.Header.Get("Accept")
		if strings.Contains(accept, "application/json") {
			return utils.JSONError(c, appErrors.New(http.StatusUnauthorized, "PASSWORD_REQUIRED", "this link is password protected", map[string]string{
				"short_code": code,
			}))
		}

		htmlForm := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Protected Link</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; background: #0f172a; color: #f8fafc; }
        .card { background: #1e293b; padding: 2.2rem; border-radius: 14px; box-shadow: 0 15px 30px rgba(0,0,0,0.5); width: 100%%; max-width: 360px; text-align: center; border: 1px solid #334155; }
        h2 { margin-top: 0; font-size: 1.35rem; font-weight: 700; color: #38bdf8; }
        p { color: #94a3b8; font-size: 0.9rem; margin-bottom: 1.5rem; line-height: 1.4; }
        input[type="password"] { width: 100%%; padding: 0.8rem; border: 1px solid #334155; border-radius: 8px; background: #0f172a; color: #fff; box-sizing: border-box; margin-bottom: 1.2rem; font-size: 1rem; outline: none; }
        input[type="password"]:focus { border-color: #38bdf8; }
        button { width: 100%%; padding: 0.8rem; background: #0284c7; color: white; border: none; border-radius: 8px; font-weight: 600; font-size: 1rem; cursor: pointer; transition: background 0.2s; }
        button:hover { background: #0369a1; }
    </style>
</head>
<body>
    <div class="card">
        <h2>🔒 Protected Link</h2>
        <p>This destination is password protected. Please enter the passkey to proceed.</p>
        <form method="POST" action="/%s">
            <input type="password" name="password" placeholder="Enter password" required autofocus />
            <button type="submit">Unlock & Continue</button>
        </form>
    </div>
</body>
</html>`, code)
		return c.HTML(http.StatusOK, htmlForm)
	}

	return c.Redirect(http.StatusFound, targetURL)
}

func (h *Handler) renderSocialPreview(c echo.Context, link *RedirectLink) error {
	title := html.EscapeString(link.Title)
	if title == "" {
		title = "Shared Link"
	}
	desc := html.EscapeString(link.Description)
	img := html.EscapeString(link.ImageURL)
	target := html.EscapeString(link.OriginalURL)

	metaHTML := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>%s</title>
    <meta property="og:title" content="%s">
    <meta property="og:description" content="%s">
    <meta property="og:image" content="%s">
    <meta property="og:type" content="website">
    <meta property="og:url" content="%s">
    <meta name="twitter:card" content="summary_large_image">
    <meta name="twitter:title" content="%s">
    <meta name="twitter:description" content="%s">
    <meta name="twitter:image" content="%s">
    <meta http-equiv="refresh" content="0; url=%s">
</head>
<body>
    <p>Redirecting to <a href="%s">%s</a>...</p>
</body>
</html>`, title, title, desc, img, target, title, desc, img, target, target, target)

	return c.HTML(http.StatusOK, metaHTML)
}

func (h *Handler) renderErrorPage(c echo.Context, code string, err error) error {
	statusCode := http.StatusNotFound
	title := "Link Not Found"
	icon := "🔍"
	desc := "The shortlink you are trying to visit does not exist, or may have been deleted."

	errMsg := err.Error()
	if strings.Contains(errMsg, "expired") {
		statusCode = http.StatusGone
		title = "Link Expired"
		icon = "⏳"
		desc = "This shortlink was set with an expiration deadline and is no longer active."
	} else if strings.Contains(errMsg, "maximum click limit") {
		statusCode = http.StatusGone
		title = "Click Quota Reached"
		icon = "🚫"
		desc = "This promotional shortlink has reached its maximum allowed clicks limit."
	} else if strings.Contains(errMsg, "deactivated") {
		statusCode = http.StatusForbidden
		title = "Link Deactivated"
		icon = "⏸️"
		desc = "This shortlink has been temporarily disabled by its creator."
	}

	htmlPage := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; background: #0f172a; color: #f8fafc; }
        .card { background: #1e293b; padding: 2.5rem; border-radius: 16px; box-shadow: 0 15px 35px rgba(0,0,0,0.5); width: 100%%; max-width: 400px; text-align: center; border: 1px solid #334155; }
        .icon { font-size: 3rem; margin-bottom: 1rem; }
        h1 { margin: 0 0 0.5rem 0; font-size: 1.5rem; font-weight: 700; color: #f87171; }
        p { color: #94a3b8; font-size: 0.95rem; line-height: 1.5; margin: 0 0 1.5rem 0; }
        .code { display: inline-block; padding: 0.25rem 0.6rem; background: #0f172a; border-radius: 6px; font-family: monospace; color: #38bdf8; font-size: 0.9rem; margin-bottom: 1.5rem; }
        a.btn { display: inline-block; padding: 0.75rem 1.5rem; background: #334155; color: #f8fafc; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 0.9rem; transition: background 0.2s; }
        a.btn:hover { background: #475569; }
    </style>
</head>
<body>
    <div class="card">
        <div class="icon">%s</div>
        <h1>%s</h1>
        <div class="code">/%s</div>
        <p>%s</p>
        <a href="/" class="btn">Go Home</a>
    </div>
</body>
</html>`, title, icon, title, code, desc)

	return c.HTML(statusCode, htmlPage)
}
