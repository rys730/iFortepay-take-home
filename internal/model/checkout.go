package model

type CheckoutItemData struct {
	ID       int32 `json:"id"`
	Quantity int32 `json:"quantity"`
}

type CheckoutRequest struct {
	// CustomerID int      `json:"customer_id"`
	Items []CheckoutItemData `json:"items"`
	// PaymentMethod string `json:"payment_method"`
}
