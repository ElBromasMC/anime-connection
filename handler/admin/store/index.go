package store

import (
	"alc/handler/util"
	"alc/view/admin/store"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET "/admin/tienda"
func (h *Handler) HandleIndexShow(c echo.Context) error {
	return util.Render(c, http.StatusOK, store.Show())
}
