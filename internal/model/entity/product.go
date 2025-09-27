package entity

import (
	"time"

	"github.com/rys730/iFortepay-take-home/db"
)

type Product struct {
	ID        int32
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Sku       string
	Name      string
	Price     float64
	Quantity  int32
}

func ProductFromDB(dbProduct *db.Product) (Product, error) {
	pricePg, err := dbProduct.Price.Float64Value()
	if err != nil {
		return Product{}, err
	}
	return Product{
		ID:        dbProduct.ID,
		CreatedAt: dbProduct.CreatedAt.Time,
		UpdatedAt: dbProduct.UpdatedAt.Time,
		DeletedAt: dbProduct.DeletedAt.Time,
		Sku:       dbProduct.Sku,
		Name:      dbProduct.Name,
		Price:     pricePg.Float64,
		Quantity:  dbProduct.Quantity,
	}, nil
}

type CheckoutItem struct {
	ProductID   int32   `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int32   `json:"quantity"`
	TotalPrice  float64 `json:"total_price"` // Price after applying promotion
}
