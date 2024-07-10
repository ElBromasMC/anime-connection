package public

import (
	"alc/config"
	"alc/handler/util"
	"alc/view/page"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleIndexShow(c echo.Context) error {
	// Query data
	items, err := h.PublicService.GetLatestItems(config.LATEST_PAGINATION)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, page.Index(items))
}

func (h *Handler) HandleNosotrosShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, page.Nosotros())
}

func (h *Handler) HandleContactoShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, page.Contacto())
}
