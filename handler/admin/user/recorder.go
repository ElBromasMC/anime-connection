package user

import (
	"alc/handler/util"
	"alc/model/auth"
	"alc/view/admin/user/recorder"
	"net/http"
	"net/mail"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) HandleRecordersShow(c echo.Context) error {
	users, err := h.AuthService.GetRecorderUsers()
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, recorder.Show(users))
}

func (h *Handler) HandleRecorderInsertion(c echo.Context) error {
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
	if err := h.AuthService.InsertRecorderUser(u, hpass); err != nil {
		return err
	}

	// Get updated users
	users, err := h.AuthService.GetRecorderUsers()
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, recorder.Table(users))
}

func (h *Handler) HandleRecorderDeletion(c echo.Context) error {
	userRawId := c.Param("userId")

	userId, err := uuid.FromString(userRawId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de usuario inválido")
	}

	// Delete recorder
	if err := h.AuthService.DeleteUser(userId); err != nil {
		return err
	}

	// Get updated users
	users, err := h.AuthService.GetRecorderUsers()
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, recorder.Table(users))
}

func (h *Handler) HandleRecorderInsertionFormShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, recorder.InsertionForm())
}

func (h *Handler) HandleRecorderDeletionFormShow(c echo.Context) error {
	userRawId := c.Param("userId")

	userId, err := uuid.FromString(userRawId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Id de usuario inválido")
	}

	user, err := h.AuthService.GetUser(userId)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, recorder.DeletionForm(user))
}
