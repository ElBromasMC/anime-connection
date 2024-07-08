package public

import (
	"alc/config"
	"alc/handler/util"
	"alc/model/auth"
	"alc/model/store"
	view "alc/view/store"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GET "/store/categories/:categorySlug/items/:itemSlug/comments?p="
func (h *Handler) HandleCommentsShow(c echo.Context) error {
	// Parsing request
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")
	page, err := strconv.Atoi(c.QueryParam("p"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Número de página no válida")
	}

	// Query data
	cat, err := h.PublicService.GetCategory(store.MangaType, categorySlug)
	if err != nil {
		return err
	}
	item, err := h.PublicService.GetItem(cat, itemSlug)
	if err != nil {
		return err
	}

	// Get comments
	comments, err := h.CommentService.GetComments(item, config.COMMENT_PAGINATION, page)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, view.Comment(item, comments, page))

}

// POST "/store/categories/:categorySlug/items/:itemSlug/comments"
func (h *Handler) HandleCommentInsertion(c echo.Context) error {
	// Parsing request
	categorySlug := c.Param("categorySlug")
	itemSlug := c.Param("itemSlug")

	var i store.ItemComment
	i.Title = c.FormValue("title")
	rating, err := strconv.Atoi(c.FormValue("rating"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Rating no válido")
	}
	if !(1 <= rating && rating <= 5) {
		return echo.NewHTTPError(http.StatusBadRequest, "Rating debe estar entre 1 y 5")
	}
	i.Rating = rating
	i.Message = c.FormValue("message")

	// Query data
	user, ok := auth.GetUser(c.Request().Context())
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "No está autenticado")
	}

	cat, err := h.PublicService.GetCategory(store.MangaType, categorySlug)
	if err != nil {
		return err
	}
	item, err := h.PublicService.GetItem(cat, itemSlug)
	if err != nil {
		return err
	}

	// Attach data
	i.Item = item
	i.CommentedBy = user

	// Insert comment
	if _, err := h.CommentService.InsertComment(i); err != nil {
		return err
	}

	// Get comments
	comments, err := h.CommentService.GetComments(item, config.COMMENT_PAGINATION, 1)
	if err != nil {
		return err
	}

	return util.Render(c, http.StatusOK, view.Comment(item, comments, 1))
}

func (h *Handler) HandleCommentUpdate(c echo.Context) error {
	return nil
}

func (h *Handler) HandleCommentDeletion(c echo.Context) error {
	return nil
}

func (h *Handler) HandleCommentUpdateFormShow(c echo.Context) error {
	return nil
}
