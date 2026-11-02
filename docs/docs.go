package docs

import (
	_ "embed"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed openapi.json
var OpenAPIJSON []byte

const swaggerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Swagger UI - Ultra-Fast Redirector & Click Collector</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
    <style>
      html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
      *, *:before, *:after { box-sizing: inherit; }
      body { margin: 0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
    <script>
      window.onload = function() {
        window.ui = SwaggerUIBundle({
          url: "/swagger/doc.json",
          dom_id: '#swagger-ui',
          deepLinking: true,
          presets: [
            SwaggerUIBundle.presets.apis,
            SwaggerUIStandalonePreset
          ],
          plugins: [
            SwaggerUIBundle.plugins.DownloadUrl
          ],
          layout: "StandaloneLayout"
        });
      };
    </script>
</body>
</html>`

// RegisterDocsEndpoint exposes the Swagger documentation at /swagger/index.html and other endpoints
func RegisterDocsEndpoint(e *echo.Echo) {
	// OpenAPI specification JSON
	serveJSON := func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/json", OpenAPIJSON)
	}
	e.GET("/swagger/doc.json", serveJSON)
	e.GET("/docs/openapi.json", serveJSON)

	// Swagger UI HTML endpoints
	serveHTML := func(c echo.Context) error {
		return c.HTML(http.StatusOK, swaggerHTML)
	}

	// Canonical requested endpoint: /swagger/index.html
	e.GET("/swagger/index.html", serveHTML)
	e.HEAD("/swagger/index.html", serveHTML)
	e.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	e.GET("/swagger/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// Also keep /docs
	e.GET("/docs", serveHTML)
	e.HEAD("/docs", serveHTML)
}
