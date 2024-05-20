package public

import (
	"alc/handler/util"
	"alc/view/page"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleIndexShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, page.Index())
}

func (h *Handler) HandleTicketShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, page.Ticket())
}
