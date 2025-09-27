package entity

import (
	"time"
)

type ProductPromotion struct {
	PromotionType string
	ID            int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
	PromotionID   int32
	ProductID     int32
	MinQuantity   int32
	FreeProductID *int32
	Discount      *float64
	FreeQuantity  *int32
	PayY          *int32
}

