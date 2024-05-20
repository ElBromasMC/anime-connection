package service

import (
	"alc/config"
	"alc/model/store"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

type Admin struct {
	Public
}

func NewAdminService(ps Public) Admin {
	return Admin{
		Public: ps,
	}
}

// Category management

func (as Admin) InsertCategory(cat store.Category) (int, error) {
	var imgId *int
	if cat.Img.Id != 0 {
		imgId = &cat.Img.Id
	}

	// Insert new category
	var id int
	if err := as.DB.QueryRow(context.Background(), `INSERT INTO store_categories (type, name, description, img_id, slug)
VALUES ($1, $2, $3, $4, $5)
RETURNING id`, cat.Type, cat.Name, cat.Description, imgId, cat.Slug).Scan(&id); err != nil {
		return 0, echo.NewHTTPError(http.StatusInternalServerError, "Error inserting new category into database")
	}
	return id, nil
}

func (as Admin) UpdateCategory(id int, uptCat store.Category) error {
	cat, err := as.GetCategoryById(id)
	if err != nil {
		return err
	}

	var imgId *int
	if cat.Img.Id != 0 {
		imgId = &cat.Img.Id
	}
	if uptCat.Img.Id != 0 {
		// Remove previous image if exists
		if cat.Img.Id != 0 {
			as.RemoveImage(cat.Img.Id)
		}
		imgId = &uptCat.Img.Id
	}

	// Update category
	if _, err := as.DB.Exec(context.Background(), `UPDATE store_categories
SET name = $1, description = $2, img_id = $3, slug = $4
WHERE id = $5`, uptCat.Name, uptCat.Description, imgId, uptCat.Slug, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error updating the category into database")
	}

	return nil
}

func (as Admin) RemoveCategory(id int) error {
	// Retrieve category
	category, err := as.GetCategoryById(id)
	if err != nil {
		return err
	}

	// Remove image of category if exists
	if category.Img.Id != 0 {
		as.RemoveImage(category.Img.Id)
	}

	// Retrive asociated items
	items, err := as.GetItems(category)
	if err != nil {
		return err
	}

	// Remove images of associated items if exists
	for _, i := range items {
		if i.Img.Id != 0 {
			as.RemoveImage(i.Img.Id)
		}
		if i.LargeImg.Id != 0 {
			as.RemoveImage(i.LargeImg.Id)
		}
	}

	if _, err := as.DB.Exec(context.Background(), `DELETE FROM store_categories
WHERE id = $1`, id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Category not found")
	}

	return nil
}

// Item management

func (as Admin) InsertItem(item store.Item) (int, error) {
	var imgId *int
	if item.Img.Id != 0 {
		imgId = &item.Img.Id
	}

	var largeimgId *int
	if item.LargeImg.Id != 0 {
		largeimgId = &item.LargeImg.Id
	}

	// Insert new item
	var id int
	if err := as.DB.QueryRow(context.Background(), `INSERT INTO store_items (category_id, name, description, long_description, img_id, largeimg_id, slug)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id`, item.Category.Id, item.Name, item.Description, item.LongDescription, imgId, largeimgId, item.Slug).Scan(&id); err != nil {
		return 0, echo.NewHTTPError(http.StatusInternalServerError, "Error inserting new item into database")
	}
	return id, nil
}

func (as Admin) UpdateItem(id int, uptItem store.Item) error {
	item, err := as.GetItemById(id)
	if err != nil {
		return err
	}

	var imgId *int
	if item.Img.Id != 0 {
		imgId = &item.Img.Id
	}
	if uptItem.Img.Id != 0 {
		// Remove previous image if exists
		if item.Img.Id != 0 {
			as.RemoveImage(item.Img.Id)
		}
		imgId = &uptItem.Img.Id
	}

	var largeimgId *int
	if item.LargeImg.Id != 0 {
		largeimgId = &item.LargeImg.Id
	}
	if uptItem.LargeImg.Id != 0 {
		// Remove previous image if exists
		if item.LargeImg.Id != 0 {
			as.RemoveImage(item.LargeImg.Id)
		}
		largeimgId = &uptItem.LargeImg.Id
	}

	// Update item
	if _, err := as.DB.Exec(context.Background(), `UPDATE store_items
SET name = $1, description = $2, long_description = $3, img_id = $4, largeimg_id = $5, slug = $6
WHERE id = $7`, uptItem.Name, uptItem.Description, uptItem.LongDescription, imgId, largeimgId, uptItem.Slug, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error updating the item into database")
	}

	return nil
}

func (as Admin) RemoveItem(id int) error {
	var imgId *int
	var largeimgId *int

	// Remove item from database
	if err := as.DB.QueryRow(context.Background(), `DELETE FROM store_items
WHERE id = $1 RETURNING img_id, largeimg_id`, id).Scan(&imgId, &largeimgId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Item not found")
	}

	// Remove attached images if exists
	if imgId != nil {
		as.RemoveImage(*imgId)
	}
	if largeimgId != nil {
		as.RemoveImage(*largeimgId)
	}

	return nil
}

// Product management

func (as Admin) InsertProduct(product store.Product) (int, error) {
	hstoreDetails := make(pgtype.Hstore, len(product.Details))
	for key, val := range product.Details {
		valCopy := val
		hstoreDetails[key] = &valCopy
	}

	var id int
	if err := as.DB.QueryRow(context.Background(), `INSERT INTO store_products (item_id, name, price, details, slug, stock)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id`, product.Item.Id, product.Name, product.Price, hstoreDetails, product.Slug, product.Stock).Scan(&id); err != nil {
		return 0, echo.NewHTTPError(http.StatusInternalServerError, "Error inserting product into database")
	}
	return id, nil
}

func (as Admin) UpdateProduct(id int, product store.Product) error {
	hstoreDetails := make(pgtype.Hstore, len(product.Details))
	for key, val := range product.Details {
		valCopy := val
		hstoreDetails[key] = &valCopy
	}

	if _, err := as.DB.Exec(context.Background(), `UPDATE store_products
SET name = $1, price = $2, details = $3, slug = $4
WHERE id = $5`, product.Name, product.Price, hstoreDetails, product.Slug, id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	return nil
}

func (as Admin) UpdateStock(id int, quantity int) error {
	sql := `UPDATE store_products SET stock = stock + $1 WHERE id = $2 AND stock + $1 >= 0`
	c, err := as.DB.Exec(context.Background(), sql, quantity, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if c.RowsAffected() != 1 {
		return echo.NewHTTPError(http.StatusNotFound, "Producto no encontrado o stock inválido")
	}
	return nil
}

func (as Admin) RemoveProduct(id int) error {
	if _, err := as.DB.Exec(context.Background(), `DELETE FROM store_products WHERE id = $1`, id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	return nil
}

// Image management

func (as Admin) InsertImage(img *multipart.FileHeader) (store.Image, error) {
	// Source
	src, err := img.Open()
	if err != nil {
		return store.Image{}, echo.NewHTTPError(http.StatusInternalServerError, "Error opening the image")
	}
	defer src.Close()

	// Destination
	dst, err := os.CreateTemp(config.IMAGES_SAVEDIR, "*"+filepath.Ext(img.Filename))
	if err != nil {
		return store.Image{}, echo.NewHTTPError(http.StatusInternalServerError, "Error creating new image")
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return store.Image{}, echo.NewHTTPError(http.StatusInternalServerError, "Error saving new image")
	}

	// Insert the image into database
	var newImg store.Image
	newImg.Filename = filepath.Base(dst.Name())
	if err := as.DB.QueryRow(context.Background(), `INSERT INTO images (filename)
VALUES ($1)
RETURNING id`, newImg.Filename).Scan(&newImg.Id); err != nil {
		return store.Image{}, echo.NewHTTPError(http.StatusInternalServerError, "Error inserting new image into database")
	}
	return newImg, nil
}

func (as Admin) RemoveImage(id int) error {
	var filename string
	// Delete from database
	if err := as.DB.QueryRow(context.Background(), `DELETE FROM images WHERE id = $1 RETURNING filename`, id).Scan(&filename); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Image not found")
	}

	// Delete from filesystem
	if err := os.Remove(path.Join(config.IMAGES_SAVEDIR, filename)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error deleting the image from filesystem")
	}
	return nil
}
