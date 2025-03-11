package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Middleware để kiểm tra token từ request header
func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		if token != "Bearer vinh" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		return next(c)
	}
}
