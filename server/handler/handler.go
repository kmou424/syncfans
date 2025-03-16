package handler

import "github.com/labstack/echo/v4"

type IHandler interface {
	Register(e *echo.Echo)
}
