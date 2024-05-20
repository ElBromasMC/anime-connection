package public

import (
	"alc/config"
	"alc/handler/util"
	"alc/model/store"
	view "alc/view/store"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GET "/store/categories/all"
func (h *Handler) HandleStoreAllShow(c echo.Context) error {
	// Query data
	cats, err := h.PublicService.GetCategories(store.StoreType)
	if err != nil {
		return err
	}
	items, err := h.PublicService.GetAllItemsLike(store.StoreType, "", 1, config.PAGINATION)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, view.Show(cats, "all", items))
}

// GET "/store/categories/:categorySlug"
func (h *Handler) HandleStoreCategoryShow(c echo.Context) error {
	// Parsing request
	categorySlug := c.Param("categorySlug")

	// Query data
	cat, err := h.PublicService.GetCategory(store.StoreType, categorySlug)
	if err != nil {
		return err
	}
	items, err := h.PublicService.GetItemsLike(cat, "", 1, config.PAGINATION)
	if err != nil {
		return err
	}
	cats, err := h.PublicService.GetCategories(store.StoreType)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, view.Show(cats, cat.Slug, items))
}

// GET "/store/categories/all/items"
func (h *Handler) HandleStoreAllItemsShow(c echo.Context) error {
	// Parsing request
	like := c.FormValue("like")
	p := c.QueryParam("p")

	page, err := strconv.Atoi(p)
	if err != nil || page < 1 {
		page = 1
	}

	// Query data
	items, err := h.PublicService.GetAllItemsLike(store.StoreType, like, page, config.PAGINATION)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, view.Category(like, page, "all", items))
}

// GET "/store/categories/:categorySlug/items"
func (h *Handler) HandleStoreCategoryItemsShow(c echo.Context) error {
	// Parsing request
	categorySlug := c.Param("categorySlug")
	like := c.FormValue("like")
	p := c.QueryParam("p")

	page, err := strconv.Atoi(p)
	if err != nil || page < 1 {
		page = 1
	}

	// Query data
	cat, err := h.PublicService.GetCategory(store.StoreType, categorySlug)
	if err != nil {
		return err
	}
	items, err := h.PublicService.GetItemsLike(cat, like, page, config.PAGINATION)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, view.Category(like, page, cat.Slug, items))
}

// GET "/store/categories/:categorySlug/items/:itemSlug"
func (h *Handler) HandleStoreItemShow(c echo.Context) error {
	// Parsing request
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")

	// Query data
	cat, err := h.PublicService.GetCategory(store.StoreType, categorySlug)
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

	return util.Render(c, http.StatusOK, view.Item(item, products))
}
