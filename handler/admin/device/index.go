package device

import (
	"alc/handler/util"
	"alc/model/auth"
	"alc/view/admin/device"
	"net/http"
	"strings"
	"unicode"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleIndexShow(c echo.Context) error {
	devices, err := h.DeviceService.GetDevices(true)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, device.Show(devices))
}

func (h *Handler) HandleInsertion(c echo.Context) error {
	// Parsing request
	serial := c.FormValue("serial")

	// Remove spaces and uppercase serial
	serial = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return unicode.ToUpper(r)
	}, serial)

	// Validate serial
	if !(12 <= len(serial) && len(serial) <= 15) {
		return echo.NewHTTPError(http.StatusBadRequest, "La serie debe tener entre 12 a 15 caracteres")
	}

	// Query data
	user, ok := auth.GetUser(c.Request().Context())
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Insert device
	h.DeviceService.InsertDevice(serial, user.Email)

	// Get updated devices
	devices, err := h.DeviceService.GetDevices(true)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, device.Table(devices))
}

func (h *Handler) HandleDeactivation(c echo.Context) error {
	return nil
}

func (h *Handler) HandleInsertionFormShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, device.InsertionForm())
}

func (h *Handler) HandleHistoryShow(c echo.Context) error {
	return nil
}

func (h *Handler) HandleDeactivationFormShow(c echo.Context) error {
	return nil
}
