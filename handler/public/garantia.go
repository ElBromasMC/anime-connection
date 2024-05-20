package public

import (
	"alc/handler/util"
	"alc/model/store"
	"alc/view/garantia"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET "/garantia"
func (h *Handler) HandleGarantiaShow(c echo.Context) error {
	cats, err := h.PublicService.GetCategories(store.GarantiaType)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, garantia.Show(cats))
}

// GET "/garantia/:slug"
func (h *Handler) HandleGarantiaCategoryShow(c echo.Context) error {
	slug := c.Param("slug")

	cat, err := h.PublicService.GetCategory(store.GarantiaType, slug)
	if err != nil {
		return err
	}

	items, err := h.PublicService.GetItems(cat)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, garantia.ShowCategory(cat, items))
}

// GET "/garantia/:categorySlug/:itemSlug"
func (h *Handler) HandleGarantiaItemShow(c echo.Context) error {
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")

	cat, err := h.PublicService.GetCategory(store.GarantiaType, categorySlug)
	if err != nil {
		return err
	}

	item, err := h.PublicService.GetItem(cat, itemSlug)
	if err != nil {
		return err
	}

	products, err := h.PublicService.GetProducts(item)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, garantia.ShowItem(item, products))
}
