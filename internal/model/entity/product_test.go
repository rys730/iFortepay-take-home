package entity_test

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rys730/iFortepay-take-home/db"
	"github.com/rys730/iFortepay-take-home/internal/model/entity"
)

func TestProductFromDB(t *testing.T) {
	var dummyPrice pgtype.Numeric
	dummyPrice.Scan("99.99")
	timeNow := time.Now()
	tests := []struct {
		name      string
		dbProduct *db.Product
		want      entity.Product
		wantErr   bool
	}{
		{
			name: "valid db product",
			dbProduct: &db.Product{
				ID:        1,
				CreatedAt: pgtype.Timestamptz{Time: timeNow, Valid: true},
				UpdatedAt: pgtype.Timestamptz{Time: timeNow, Valid: true},
				DeletedAt: pgtype.Timestamptz{Time: time.Time{}, Valid: false},
				Sku:       "SKU123",
				Name:      "Test Product",
				Price:     dummyPrice,
				Quantity:  10,
			},
			want: entity.Product{
				ID:        1,
				CreatedAt: timeNow,
				UpdatedAt: timeNow,
				DeletedAt: time.Time{},
				Sku:       "SKU123",
				Name:      "Test Product",
				Price:     99.99,
				Quantity:  10,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := entity.ProductFromDB(tt.dbProduct)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ProductFromDB() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ProductFromDB() succeeded unexpectedly")
			}
			if got.ID != tt.want.ID ||
				got.Sku != tt.want.Sku ||
				got.Name != tt.want.Name ||
				got.Quantity != tt.want.Quantity {
				t.Errorf("Product mismatch: got %+v, want %+v", got, tt.want)
			}

		})
	}
}
