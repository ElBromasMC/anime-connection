package store

import (
	"alc/handler/util"
	"alc/model/store"
	"alc/view/admin/store/category"
	"net/http"

	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
)

// GET "/admin/tienda/type/:typeSlug/categories"
func (h *Handler) HandleCategoriesShow(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}

	cats, err := h.AdminService.GetCategories(t)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, category.Show(t, cats))
}

// POST "/admin/tienda/type/:typeSlug/categories"
func (h *Handler) HandleCategoryInsertion(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")

	var cat store.Category
	cat.Name = c.FormValue("name")
	cat.Description = c.FormValue("description")
	img, imgErr := c.FormFile("img")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}

	// Attach data
	cat.Type = t
	cat.Slug = slug.Make(cat.Name)

	// Insert and attach image if present in request
	if imgErr == nil {
		newImg, err := h.AdminService.InsertImage(img)
		if err != nil {
			return err
		}
		cat.Img = newImg
	}

	// Insert it into database
	if _, err := h.AdminService.InsertCategory(cat); err != nil {
		return err
	}

	// Get updated categories
	cats, err := h.AdminService.GetCategories(t)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, category.Table(t, cats))
}

// PUT "/admin/tienda/type/:typeSlug/categories/:categorySlug"
func (h *Handler) HandleCategoryUpdate(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	categorySlug := c.Param("categorySlug")

	var cat store.Category
	cat.Name = c.FormValue("name")
	cat.Description = c.FormValue("description")
	img, imgErr := c.FormFile("img")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}
	oldCat, err := h.AdminService.GetCategory(t, categorySlug)
	if err != nil {
		return err
	}

	// Attach data
	cat.Type = t
	cat.Slug = slug.Make(cat.Name)

	// Insert and attach image if present in request
	if imgErr == nil {
		newImg, err := h.AdminService.InsertImage(img)
		if err != nil {
			return err
		}
		cat.Img = newImg
	}

	// Update category
	if err := h.AdminService.UpdateCategory(oldCat.Id, cat); err != nil {
		return err
	}

	// Get updated categories
	cats, err := h.AdminService.GetCategories(t)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, category.Table(t, cats))
}

// DELETE "/admin/tienda/type/:typeSlug/categories/:categorySlug"
func (h *Handler) HandleCategoryDeletion(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")
	categorySlug := c.Param("categorySlug")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}
	oldCat, err := h.AdminService.GetCategory(t, categorySlug)
	if err != nil {
		return err
	}

	// Remove category
	if err := h.AdminService.RemoveCategory(oldCat.Id); err != nil {
		return err
	}

	// Get updated categories
	cats, err := h.AdminService.GetCategories(t)
	if err != nil {
		return err
	}
	return util.Render(c, http.StatusOK, category.Table(t, cats))
}

// GET "/admin/tienda/type/:typeSlug/categories/insert"
func (h *Handler) HandleCategoryInsertionFormShow(c echo.Context) error {
	// Parsing request
	typeSlug := c.Param("typeSlug")

	// Query data
	t, err := h.AdminService.GetType(typeSlug)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, category.InsertionForm(t))
}

// GET "/admin/tienda/type/:typeSlug/categories/:categorySlug/update"
func (h *Handler) HandleCategoryUpdateFormShow(c echo.Context) error {
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

	return util.Render(c, http.StatusOK, category.UpdateForm(cat))
}

// GET "/admin/tienda/type/:typeSlug/categories/:categorySlug/delete"
func (h *Handler) HandleCategoryDeletionFormShow(c echo.Context) error {
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

	return util.Render(c, http.StatusOK, category.DeletionForm(cat))
}
