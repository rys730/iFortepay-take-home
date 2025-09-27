package promotion

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/rys730/iFortepay-take-home/db"
	"github.com/rys730/iFortepay-take-home/internal/model/entity"
	"github.com/stretchr/testify/require"
)

func TestPromotionRepository_GetProductPromotionsByProductID(t *testing.T) {
	ctx := context.Background()

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Failed to create mock pool: %v", err)
	}
	defer mock.Close()

	repo := &PromotionRepository{db: mock}

	productID := int32(1)
	var num pgtype.Numeric
	num.Scan("0.10")
	insertedPromotions := []db.GetProductPromotionsByProductIDRow{
		{
			PromotionType: "FREE_ITEM",
			ID:            1,
			CreatedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
			UpdatedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
			DeletedAt:     pgtype.Timestamptz{Time: time.Time{}, Valid: false},
			PromotionID:   1,
			ProductID:     1,
			MinQuantity:   2,
			FreeProductID: &[]int32{2}[0],
			Discount:      num,
			FreeQuantity:  &[]int32{1}[0],
			PayY:          &[]int32{1}[0],
		},
	}

	var expectedPromotions []entity.ProductPromotion
	pricePg, err := insertedPromotions[0].Discount.Float64Value()
	require.NoError(t, err)

	var discount *float64
	if pricePg.Valid {
		discount = &pricePg.Float64
	}
	pp := entity.ProductPromotion{
		PromotionType: insertedPromotions[0].PromotionType,
		ID:            insertedPromotions[0].ID,
		CreatedAt:     insertedPromotions[0].CreatedAt.Time,
		UpdatedAt:     insertedPromotions[0].UpdatedAt.Time,
		DeletedAt:     insertedPromotions[0].DeletedAt.Time,
		PromotionID:   insertedPromotions[0].PromotionID,
		ProductID:     insertedPromotions[0].ProductID,
		MinQuantity:   insertedPromotions[0].MinQuantity,
		FreeProductID: insertedPromotions[0].FreeProductID,
		Discount:      discount,
		FreeQuantity:  insertedPromotions[0].FreeQuantity,
		PayY:          insertedPromotions[0].PayY,
	}
	expectedPromotions = append(expectedPromotions, pp)

	rows := pgxmock.NewRows([]string{"promotion_type", "id", "created_at", "updated_at", "deleted_at", "promotion_id", "product_id", "min_quantity", "free_product_id", "discount", "free_quantity", "pay_y"}).
		AddRow(insertedPromotions[0].PromotionType, insertedPromotions[0].ID, insertedPromotions[0].CreatedAt, insertedPromotions[0].UpdatedAt, insertedPromotions[0].DeletedAt, insertedPromotions[0].PromotionID, insertedPromotions[0].ProductID, insertedPromotions[0].MinQuantity, insertedPromotions[0].FreeProductID, insertedPromotions[0].Discount, insertedPromotions[0].FreeQuantity, insertedPromotions[0].PayY)

	mock.ExpectQuery(regexp.QuoteMeta(`select p.promotion_type, pp.id, pp.created_at, pp.updated_at, pp.deleted_at, pp.promotion_id, pp.product_id, pp.min_quantity, pp.free_product_id, pp.discount, pp.free_quantity, pp.pay_y from promotions p join product_promotions pp on pp.promotion_id = p.id where p.start_date <= CURRENT_DATE and p.end_date >= CURRENT_DATE and pp.product_id = $1 and pp.deleted_at is null and p.deleted_at is null`)).WithArgs(productID).WillReturnRows(rows)

	promotions, err := repo.GetProductPromotionsByProductID(ctx, productID)
	require.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if !reflect.DeepEqual(promotions, expectedPromotions) {
		t.Errorf("Promotions is not as expected: \nExpected: %v\nGot: %v", insertedPromotions, promotions)
	}
}

func TestPromotionRepository_GetProductPromotionsByProductID_ProductNotFound(t *testing.T) {
	ctx := context.Background()

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Failed to create mock pool: %v", err)
	}
	defer mock.Close()

	repo := &PromotionRepository{db: mock}

	productID := int32(1)

	mock.ExpectQuery(regexp.QuoteMeta(`select p.promotion_type, pp.id, pp.created_at, pp.updated_at, pp.deleted_at, pp.promotion_id, pp.product_id, pp.min_quantity, pp.free_product_id, pp.discount, pp.free_quantity, pp.pay_y from promotions p join product_promotions pp on pp.promotion_id = p.id where p.start_date <= CURRENT_DATE and p.end_date >= CURRENT_DATE and pp.product_id = $1 and pp.deleted_at is null and p.deleted_at is null`)).WithArgs(productID).WillReturnError(errors.New("product not found"))

	_, err = repo.GetProductPromotionsByProductID(ctx, productID)
	require.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if err.Error() != "product not found" {
		t.Errorf("Error is not as expected: \nExpected: %v\nGot: %v", "product not found", err.Error())
	}
}
