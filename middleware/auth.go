package middleware

import (
	"alc/model/auth"
	"alc/service"
	"context"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Auth(us service.Auth) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get(auth.SessionName, c)
			if err != nil {
				c.Logger().Debug("Error getting user session: ", err)
				return next(c)
			}

			// Retrieve session
			s, ok := sess.Values["session"].([]byte)
			if !ok {
				c.Logger().Debug("Invalid session from request")
				return next(c)
			}

			// Get UUID
			sessionID, err := uuid.FromBytes(s)
			if err != nil {
				c.Logger().Debug("Invalid session from request: ", err)
				return next(c)
			}

			// Get user
			u, err := us.GetUserBySession(sessionID)
			if err != nil {
				c.Logger().Debug("Unauthorized: ", err)
				return next(c)
			}

			// Attach user to request context
			ctx := context.WithValue(c.Request().Context(), auth.AuthKey{}, u)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func Admin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		u, ok := auth.GetUser(c.Request().Context())
		if !ok {
			return c.Redirect(http.StatusFound, "/login?to=/admin")
		}
		if u.Role != auth.AdminRole && u.Role != auth.RecorderRole {
			return c.Redirect(http.StatusFound, "/login?to=/admin")
		}
		return next(c)
	}
}

func RoleAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		u, ok := auth.GetUser(c.Request().Context())
		if !ok {
			return c.Redirect(http.StatusFound, "/login?to=/admin")
		}
		if u.Role != auth.AdminRole {
			return c.Redirect(http.StatusFound, "/login?to=/admin")
		}
		return next(c)
	}
}
