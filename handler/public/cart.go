package public

import (
	"alc/handler/util"
	"alc/model/cart"
	"alc/model/store"
	view "alc/view/cart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// POST "/cart"
func (h *Handler) HandleNewCartItem(c echo.Context) error {
	// Bind
	var i cart.ItemRequest
	id, err := strconv.Atoi(c.FormValue("Id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id")
	}
	i.ProductId = id
	quantity, err := strconv.Atoi(c.FormValue("Quantity"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Quantity")
	}
	i.Quantity = quantity
	serie := c.FormValue("Details[Serie]")
	if len(serie) != 0 {
		i.Details = make(map[string]string)
		i.Details["Serie"] = serie
	}

	// Get and validate the new item
	item, err := i.ToItem(h.PublicService)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	if err := item.IsValid(); err != nil {
		if err.Error() == "quantity exceeds current stock" {
			return echo.NewHTTPError(http.StatusBadRequest, "La cantidad seleccionada supera al stock")
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	// Get cart items
	items := cart.GetItems(c.Request().Context())

	// Insert the new item
	found := false
	for n, i := range items {
		if i.Product.Item.Id != item.Product.Item.Id {
			continue
		}
		if i.Product.Item.Category.Type != store.StoreType {
			if !strings.EqualFold(i.Details["Serie"], item.Details["Serie"]) {
				continue
			}
			found = true
			items[n] = item
			continue
		}
		if i.Product.Id != item.Product.Id {
			continue
		}
		found = true
		items[n].Quantity += item.Quantity
	}
	if !found {
		items = append(items, item)
	}

	// Validate new items
	for _, i := range items {
		if err := i.IsValid(); err != nil {
			if err.Error() == "quantity exceeds current stock" {
				return echo.NewHTTPError(http.StatusBadRequest, "La cantidad seleccionada supera al stock")
			} else {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		}
	}

	// Update cart items
	itemsRequest := make([]cart.ItemRequest, 0, len(items))
	for _, i := range items {
		itemsRequest = append(itemsRequest, i.ToRequest())
	}
	sess, _ := session.Get(cart.SessionName, c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   24 * 60 * 60,
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
	}
	sess.Values["items"] = itemsRequest
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		c.Logger().Debug("Error saving cart session: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error al guardar el item")
	}

	return util.Render(c, http.StatusOK, view.Show(items))
}

// DELETE "/cart"
func (h *Handler) HandleRemoveCartItem(c echo.Context) error {
	// Bind
	id, err := strconv.Atoi(c.FormValue("Id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id")
	}

	// Get cart items
	prevItems := cart.GetItems(c.Request().Context())

	// Remove cart item
	items := make([]cart.Item, 0, len(prevItems))
	for n, i := range prevItems {
		if n != id {
			items = append(items, i)
		}
	}

	// Update cart items
	itemsRequest := make([]cart.ItemRequest, 0, len(items))
	for _, i := range items {
		itemsRequest = append(itemsRequest, i.ToRequest())
	}
	sess, _ := session.Get(cart.SessionName, c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   24 * 60 * 60,
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
	}
	sess.Values["items"] = itemsRequest
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		c.Logger().Debug("Error saving cart session: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error al eliminar el item")
	}

	return util.Render(c, http.StatusOK, view.Show(items))
}
