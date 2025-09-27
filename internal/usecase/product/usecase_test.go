package product

import (
	"reflect"
	"testing"
)

func TestNewProductUsecase(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockPromoRepo := new(MockPromotionRepository)

	pu := NewProductUsecase(mockProductRepo, mockPromoRepo)

	if pu == nil {
		t.Errorf("NewProductUsecase returned nil")
	}

	expectedType := "*product.ProductUsecase"
	actualType := reflect.TypeOf(pu).String()

	if actualType != expectedType {
		t.Errorf("NewProductUsecase returned incorrect type: expected %s, got %s", expectedType, actualType)
	}
}