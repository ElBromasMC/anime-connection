package store

import (
	"alc/handler/util"
	"alc/model/store"
	"alc/view/admin/store/item"
	"net/http"

	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
)

// GET "/admin/tienda/type/:typeSlug/categories/:categorySlug/items"
func (h *Handler) HandleItemsShow(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	categorySlug := c.Param("categorySlug")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}
	cat, err := h.AdminService.GetCategory(t, categorySlug)
	if err != nil {
		return err
	}

	items, err := h.AdminService.GetItems(cat)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, item.Show(cat, items))
}

// POST "/admin/tienda/type/:typeSlug/categories/:categorySlug/items"
func (h *Handler) HandleItemInsertion(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	categorySlug := c.Param("categorySlug")

	var i store.Item
	i.Name = c.FormValue("name")
	i.Description = c.FormValue("description")
	i.LongDescription = c.FormValue("longDescription")
	img, imgErr := c.FormFile("img")
	largeImg, largeImgErr := c.FormFile("largeImg")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}
	cat, err := h.AdminService.GetCategory(t, categorySlug)
	if err != nil {
		return err
	}

	// Attach data
	i.Category = cat
	i.Slug = slug.Make(i.Name)

	// Insert and attach images if they are present in request
	if imgErr == nil {
		newImg, err := h.AdminService.InsertImage(img)
		if err != nil {
			return err
		}
		i.Img = newImg
	}

	if largeImgErr == nil {
		newLargeImg, err := h.AdminService.InsertImage(largeImg)
		if err != nil {
			return err
		}
		i.LargeImg = newLargeImg
	}

	// Insert it into database
	if _, err := h.AdminService.InsertItem(i); err != nil {
		return err
	}

	// Get updated items
	items, err := h.AdminService.GetItems(cat)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, item.Table(cat, items))
}

// PUT "/admin/tienda/type/:typeSlug/categories/:categorySlug/items/:itemSlug"
func (h *Handler) HandleItemUpdate(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")

	var i store.Item
	i.Name = c.FormValue("name")
	i.Description = c.FormValue("description")
	img, imgErr := c.FormFile("img")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}
	cat, err := h.AdminService.GetCategory(t, categorySlug)
	if err != nil {
		return err
	}
	oldItem, err := h.AdminService.GetItem(cat, itemSlug)
	if err != nil {
		return err
	}

	// Attach data
	i.Category = cat
	i.Slug = slug.Make(i.Name)

	// Insert and attach image if present in request
	if imgErr == nil {
		newImg, err := h.AdminService.InsertImage(img)
		if err != nil {
			return err
		}
		i.Img = newImg
	}

	// Update item
	if err := h.AdminService.UpdateItem(oldItem.Id, i); err != nil {
		return err
	}

	// Get updated items
	items, err := h.AdminService.GetItems(cat)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, item.Table(cat, items))
}

// DELETE "/admin/tienda/type/:typeSlug/categories/:categorySlug/items/:itemSlug"
func (h *Handler) HandleItemDeletion(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}
	cat, err := h.AdminService.GetCategory(t, categorySlug)
	if err != nil {
		return err
	}
	i, err := h.AdminService.GetItem(cat, itemSlug)
	if err != nil {
		return err
	}

	// Remove item
	if err := h.AdminService.RemoveItem(i.Id); err != nil {
		return err
	}

	// Get updated items
	items, err := h.AdminService.GetItems(cat)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, item.Table(cat, items))
}

// GET "/admin/tienda/type/:typeSlug/categories/:categorySlug/items/insert"
func (h *Handler) HandleItemInsertionFormShow(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	categorySlug := c.Param("categorySlug")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}
	cat, err := h.AdminService.GetCategory(t, categorySlug)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, item.InsertionForm(cat))
}

// GET "/admin/tienda/type/:typeSlug/categories/:categorySlug/items/:itemSlug/update"
func (h *Handler) HandleItemUpdateFormShow(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}
	cat, err := h.AdminService.GetCategory(t, categorySlug)
	if err != nil {
		return err
	}
	i, err := h.AdminService.GetItem(cat, itemSlug)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, item.UpdateForm(i))
}

// GET "/admin/tienda/type/:typeSlug/categories/:categorySlug/items/:itemSlug/delete"
func (h *Handler) HandleItemDeletionFormShow(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}
	cat, err := h.AdminService.GetCategory(t, categorySlug)
	if err != nil {
		return err
	}
	i, err := h.AdminService.GetItem(cat, itemSlug)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, item.DeletionForm(i))
}
