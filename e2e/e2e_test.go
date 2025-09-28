package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/rys730/iFortepay-take-home/internal/model/dto"
	"github.com/rys730/iFortepay-take-home/internal/model/entity"
	"github.com/stretchr/testify/assert"
)

// TestE2E_Checkout still need improvement to make it more robust
// ex: starting its own server, creating test db, run migration
func TestE2E_Checkout(t *testing.T) {
	tests := []struct {
		name    string
		payload dto.CheckoutRequest
		want    dto.CheckoutResponse
		wantErr string
	}{
		{
			name: "Checkout 1 Macbook Pro Get Free 1 Raspberry Pi B",
			payload: dto.CheckoutRequest{
				Items: []dto.CheckoutItemData{
					{ID: 2, Quantity: 1},
				},
			},
			want: dto.CheckoutResponse{
				Message:    "Checkout successful",
				TotalPrice: "$5399.99",
				Items: []entity.CheckoutItem{
					{ProductID: 2, ProductName: "MacBook Pro", Quantity: 1, TotalPrice: 5399.99},
					{ProductID: 4, ProductName: "Raspberry Pi B", Quantity: 1, TotalPrice: 0.00},
				},
			},
		},
		{
			name: "Checkout 3 Google Homes For The Price Of 2",
			payload: dto.CheckoutRequest{
				Items: []dto.CheckoutItemData{
					{ID: 1, Quantity: 3},
				},
			},
			want: dto.CheckoutResponse{
				Message:    "Checkout successful",
				TotalPrice: "$99.98",
				Items: []entity.CheckoutItem{
					{ProductID: 1, ProductName: "Google Home", Quantity: 3, TotalPrice: 99.98},
				},
			},
		},
		{
			name: "Checkout 3 or more Alexa Speakers Get 10% Discount on all Alexa Speakers",
			payload: dto.CheckoutRequest{
				Items: []dto.CheckoutItemData{
					{ID: 3, Quantity: 3},
				},
			},
			want: dto.CheckoutResponse{
				Message:    "Checkout successful",
				TotalPrice: "$295.65",
				Items: []entity.CheckoutItem{
					{ProductID: 3, ProductName: "Alexa Speaker", Quantity: 3, TotalPrice: 295.65},
				},
			},
		},
		{
			name: "Checkout product with no promotion",
			payload: dto.CheckoutRequest{
				Items: []dto.CheckoutItemData{
					{ID: 1, Quantity: 1},
				},
			},
			want: dto.CheckoutResponse{
				Message:    "Checkout successful",
				TotalPrice: "$49.99",
				Items: []entity.CheckoutItem{
					{ProductID: 1, ProductName: "Google Home", Quantity: 1, TotalPrice: 49.99},
				},
			},
		},
		{
			name: "Checkout 2 Macbook Pro Get Free 2 Raspberry Pi B but only 1 Raspberry Pi B in stock",
			payload: dto.CheckoutRequest{
				Items: []dto.CheckoutItemData{
					{ID: 2, Quantity: 2},
				},
			},
			want: dto.CheckoutResponse{
				Message:    "Checkout successful",
				TotalPrice: "$10799.98",
				Items: []entity.CheckoutItem{
					{ProductID: 2, ProductName: "MacBook Pro", Quantity: 2, TotalPrice: 10799.98},
					{ProductID: 4, ProductName: "Raspberry Pi B", Quantity: 1, TotalPrice: 0.00},
				},
			},
		},
		{
			name: "Checkout multiple products with multiple promotions",
			payload: dto.CheckoutRequest{
				Items: []dto.CheckoutItemData{
					{ID: 1, Quantity: 3},
					{ID: 3, Quantity: 4},
				},
			},
			want: dto.CheckoutResponse{
				Message:    "Checkout successful",
				TotalPrice: "$494.18",
				Items: []entity.CheckoutItem{
					{ProductID: 1, ProductName: "Google Home", Quantity: 3, TotalPrice: 99.98},
					{ProductID: 3, ProductName: "Alexa Speaker", Quantity: 4, TotalPrice: 394.20},
				},
			},
		},
		{
			name: "Checkout product without stock",
			payload: dto.CheckoutRequest{
				Items: []dto.CheckoutItemData{
					{ID: 4, Quantity: 1},
				},
			},
			want:    dto.CheckoutResponse{},
			wantErr: "product with id 4 not found or out of stock",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			resp, err := http.Post("http://localhost:8085/api/checkout", "application/json", bytes.NewBuffer(jsonBody))
			assert.NoError(t, err)
			defer resp.Body.Close()

			if tt.wantErr == "" {
				var got dto.CheckoutResponse
				err = json.NewDecoder(resp.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				var got struct {
					Message string `json:"message"`
				}
				err = json.NewDecoder(resp.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantErr, got.Message)

			}
		})
	}
}
