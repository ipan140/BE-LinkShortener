package docs

import (
	_ "embed"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed openapi.json
var OpenAPIJSON []byte

// RegisterDocsEndpoint exposes the OpenAPI documentation
func RegisterDocsEndpoint(e *echo.Echo) {
	e.GET("/docs/openapi.json", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/json", OpenAPIJSON)
	})

	e.GET("/docs", func(c echo.Context) error {
		html := `<!DOCTYPE html>
<html>
<head>
    <title>API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.onload = function() {
        SwaggerUIBundle({
          url: "/docs/openapi.json",
          dom_id: '#swagger-ui',
        });
      };
    </script>
</body>
</html>`
		return c.HTML(http.StatusOK, html)
	})
}
