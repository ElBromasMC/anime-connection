package admin

import (
	"alc/handler/util"
	"alc/view/admin"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET "/admin"
func (h *Handler) HandleIndexShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, admin.Show())
}
