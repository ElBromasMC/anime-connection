package util

import (
	"alc/view/component"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler(err error, c echo.Context) {
	switch v := err.(type) {
	case *echo.HTTPError:
		msg, ok := v.Message.(string)
		if !ok {
			msg = "Error desconocido"
		}
		switch v.Code {
		case http.StatusInternalServerError:
			if err := Render(c, v.Code, component.ErrorMessage(msg)); err != nil {
				c.Logger().Error(err)
			}
		default:
			if err := Render(c, v.Code, component.ErrorMessage(msg)); err != nil {
				c.Logger().Error(err)
			}
		}
	default:
		if err := c.String(http.StatusInternalServerError, "500 - Internal Server Error"); err != nil {
			c.Logger().Error(err)
		}
	}
}
