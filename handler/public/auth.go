package public

import (
	"alc/handler/util"
	"alc/model/auth"
	view "alc/view/user"
	"net/http"
	"net/mail"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Signup
func (h *Handler) HandleSignupShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, view.SignupShow())
}

func (h *Handler) HandleSignup(c echo.Context) error {
	// Bind
	var u auth.User
	if err := c.Bind(&u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Formato inválido")
	}

	// Trim name and email
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Name = strings.TrimSpace(u.Name)

	// Validate
	address, err := mail.ParseAddress(u.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Email inválido")
	}
	u.Email = address.Address

	if len(u.Password) < 8 {
		return echo.NewHTTPError(http.StatusBadRequest, "Contraseña muy corta (mínimo 8 caracteres)")
	}

	if len(u.Password) > 72 {
		return echo.NewHTTPError(http.StatusBadRequest, "Contraseña muy larga (máximo 72 caracteres)")
	}

	if !(0 < len(u.Name) && len(u.Name) <= 200) {
		return echo.NewHTTPError(http.StatusBadRequest, "Nombre inválido")
	}

	// Hash password
	hpass, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error desconocido")
	}

	// Save user
	if err := h.AuthService.InsertUser(u, hpass); err != nil {
		return err
	}

	_, ok := c.Request().Header[http.CanonicalHeaderKey("HX-Request")]
	if !ok {
		return c.Redirect(http.StatusFound, "/login")
	}
	c.Response().Header().Set("HX-Redirect", "/login")
	return c.NoContent(http.StatusOK)
}

// Login
func (h *Handler) HandleLoginShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, view.LoginShow(c.QueryParam("to")))
}

func (h *Handler) HandleLogin(c echo.Context) error {
	// Bind
	var u auth.User
	if err := c.Bind(&u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Formato inválido")
	}

	// Trim email
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	// Verify credentials
	id, hpass, err := h.AuthService.GetUserIdAndHpassByEmail(u.Email)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword(hpass, []byte(u.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Contraseña incorrecta")
	}

	// Start new session
	s, err := h.AuthService.InsertSession(id)
	if err != nil {
		return err
	}

	// Write a cookie
	sess, _ := session.Get(auth.SessionName, c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   30 * 24 * 60 * 60,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	sess.Values["session"] = s.Bytes()
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		c.Logger().Debug("Error saving auth session: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error al guardar la sesión")
	}

	if to := c.QueryParam("to"); len(to) > 0 {
		_, ok := c.Request().Header[http.CanonicalHeaderKey("HX-Request")]
		if !ok {
			return c.Redirect(http.StatusFound, to)
		}
		c.Response().Header().Set("HX-Redirect", to)
		return c.NoContent(http.StatusOK)
	}

	_, ok := c.Request().Header[http.CanonicalHeaderKey("HX-Request")]
	if !ok {
		return c.Redirect(http.StatusFound, "/")
	}
	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
}

// Logout
func (h *Handler) HandleLogout(c echo.Context) error {
	// Get session and defer removal
	sess, _ := session.Get(auth.SessionName, c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	defer sess.Save(c.Request(), c.Response())

	// Retrieve session
	s, ok := sess.Values["session"].([]byte)
	if !ok {
		return c.Redirect(http.StatusFound, "/")
	}

	// Get UUID
	sessionID, err := uuid.FromBytes(s)
	if err != nil {
		return c.Redirect(http.StatusFound, "/")
	}

	// Delete session from database
	if err := h.AuthService.DeleteSession(sessionID); err != nil {
		return c.Redirect(http.StatusFound, "/")
	}

	return c.Redirect(http.StatusFound, "/")
}
