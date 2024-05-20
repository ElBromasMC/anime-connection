package cart

import (
	"alc/model/store"
	"alc/service"
	"context"
	"errors"
)

const SessionName = "cart"

type ItemsKey struct{}

type Item struct {
	Product  store.Product
	Quantity int
	Details  map[string]string
}

type ItemRequest struct {
	ProductId int
	Quantity  int
	Details   map[string]string
}

func (item Item) IsValid() error {
	if item.Product.Item.Category.Type == store.GarantiaType {
		if item.Quantity != 1 {
			return errors.New("invalid quantity for warranty")
		}
		serie, ok := item.Details["Serie"]
		if !ok {
			return errors.New("missing 'Serie' for warranty")
		}
		if !(12 <= len(serie) && len(serie) <= 15) {
			return errors.New("invalid 'Serie' for warranty")
		}
	} else {
		if item.Quantity < 1 {
			return errors.New("invalid quantity for store item")
		}
		if item.Product.Stock != nil {
			if item.Quantity > *item.Product.Stock {
				return errors.New("quantity exceeds current stock")
			}
		}
	}
	return nil
}

func (i Item) ToRequest() ItemRequest {
	return ItemRequest{
		ProductId: i.Product.Id,
		Quantity:  i.Quantity,
		Details:   i.Details,
	}
}

func (i ItemRequest) ToItem(ps service.Public) (Item, error) {
	product, err := ps.GetProductById(i.ProductId)
	if err != nil {
		return Item{}, err
	}
	item := Item{
		Product:  product,
		Quantity: i.Quantity,
		Details:  i.Details,
	}
	return item, nil
}

func GetItems(ctx context.Context) []Item {
	items, _ := ctx.Value(ItemsKey{}).([]Item)
	return items
}
