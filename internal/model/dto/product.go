package dto

import "github.com/rys730/iFortepay-take-home/internal/model/entity"

type CheckoutItemData struct {
	ID       int32 `json:"id"`
	Quantity int32 `json:"quantity"`
}

type CheckoutRequest struct {
	// CustomerID int      `json:"customer_id"`
	Items []CheckoutItemData `json:"items"`
	// PaymentMethod string `json:"payment_method"`
}

type CheckoutResponse struct {
	Message    string                `json:"message"`
	TotalPrice string                `json:"total_price"`
	Items      []entity.CheckoutItem `json:"items"`
}
