//go:build dev

package main

import (
	"github.com/labstack/echo/v4"
)

func static(e *echo.Echo) *echo.Route {
	return e.Static("/static", "assets/static")
}
