package checkout

import (
	"alc/model/auth"
	"alc/model/store"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Status string

const (
	Pendiente    Status = "PENDIENTE"
	EnProceso    Status = "EN PROCESO"
	PorConfirmar Status = "POR CONFIRMAR"
	Entregado    Status = "ENTREGADO"
	Cancelado    Status = "CANCELADO"
)

type Order struct {
	Id            uuid.UUID
	PurchaseOrder int
	Email         string
	Phone         string
	Name          string
	Address       string
	City          string
	PostalCode    string
	AssignedTo    auth.User
	CreatedAt     time.Time
}

type OrderProduct struct {
	Id              int
	Order           Order
	Quantity        int
	Details         map[string]string
	Product         store.Product
	ProductType     store.Type
	ProductCategory string
	ProductItem     string
	ProductName     string
	ProductPrice    int
	ProductDetails  map[string]string
	Status          Status
	UpdatedAt       time.Time
}
