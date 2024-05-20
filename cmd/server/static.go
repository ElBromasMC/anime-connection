//go:build !dev

package main

import (
	"alc/assets"

	"github.com/labstack/echo/v4"
)

func static(e *echo.Echo) *echo.Route {
	// return e.GET("/static/*", echo.WrapHandler(http.FileServer(http.FS(assets.Assets))))
	return e.StaticFS("/static", echo.MustSubFS(assets.Assets, "static"))
}
