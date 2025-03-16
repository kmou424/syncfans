package middleware

import (
	"github.com/kmou424/syncfans/internal/conf"
	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(401, "Authorization header is missing")
		}

		secret := conf.GetServerConfig().Config.Secret
		if secret == "" {
			return next(c)
		}

		if authHeader != secret {
			return c.JSON(401, "Invalid secret")
		}

		return next(c)
	}
}
