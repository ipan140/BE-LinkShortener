package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
	appErrors "link-shortener/common/errors"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func JSONSuccess(c echo.Context, status int, message string, data any) error {
	return c.JSON(status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func JSONError(c echo.Context, err error) error {
	if appErr, ok := err.(*appErrors.AppError); ok {
		return c.JSON(appErr.StatusCode, APIResponse{
			Success: false,
			Error:   appErr,
		})
	}

	return c.JSON(http.StatusInternalServerError, APIResponse{
		Success: false,
		Error: map[string]string{
			"code":    "INTERNAL_ERROR",
			"message": err.Error(),
		},
	})
}
