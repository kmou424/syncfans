package middleware

import (
	"github.com/kmou424/syncfans/internal/caused"
	"github.com/labstack/echo/v4"
)

func Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer caused.Recover(false)
		return next(c)
	}
}
