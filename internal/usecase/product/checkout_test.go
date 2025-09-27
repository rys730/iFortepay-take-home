package product

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/rys730/iFortepay-take-home/internal/model/dto"
	"github.com/rys730/iFortepay-take-home/internal/model/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetProductByID(ctx context.Context, id int32) (entity.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Product), args.Error(1)
}

func (m *MockProductRepository) UpdateProductQuantityByID(ctx context.Context, id int32, quantity int32) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

type MockPromotionRepository struct {
	mock.Mock
}

func (m *MockPromotionRepository) GetProductPromotionsByProductID(ctx context.Context, productID int32) ([]entity.ProductPromotion, error) {
	args := m.Called(ctx, productID)
	return args.Get(0).([]entity.ProductPromotion), args.Error(1)
}

func TestProductUsecase_Checkout(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockPromoRepo := new(MockPromotionRepository)

	pu := &ProductUsecase{
		productRepo: mockProductRepo,
		promoRepo:   mockPromoRepo,
	}

	ctx := context.Background()

	req := dto.CheckoutRequest{
		Items: []dto.CheckoutItemData{
			{
				ID:       1,
				Quantity: 2,
			},
		},
	}

	product := entity.Product{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Sku:       "120P90",
		Name:      "Google Home",
		Price:     49.99,
		Quantity:  10,
	}

	promotions := []entity.ProductPromotion{}

	mockProductRepo.On("GetProductByID", ctx, int32(1)).Return(product, nil).Once()
	mockPromoRepo.On("GetProductPromotionsByProductID", ctx, int32(1)).Return(promotions, nil).Once()
	mockProductRepo.On("UpdateProductQuantityByID", ctx, int32(1), int32(2)).Return(nil).Once()

	expectedResponse := dto.CheckoutResponse{
		Message:    "Checkout successful",
		TotalPrice: "$99.98",
		Items: []entity.CheckoutItem{
			{
				ProductID:   1,
				ProductName: "Google Home",
				Quantity:    2,
				TotalPrice:  99.98,
			},
		},
	}

	response, err := pu.Checkout(ctx, req)

	mockProductRepo.AssertExpectations(t)
	mockPromoRepo.AssertExpectations(t)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("Response is not as expected: \nExpected: %v\nGot: %v", expectedResponse, response)
	}
}

func TestProductUsecase_Checkout_ProductNotFound(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockPromoRepo := new(MockPromotionRepository)

	pu := &ProductUsecase{
		productRepo: mockProductRepo,
		promoRepo:   mockPromoRepo,
	}

	ctx := context.Background()

	req := dto.CheckoutRequest{
		Items: []dto.CheckoutItemData{
			{
				ID:       50,
				Quantity: 2,
			},
		},
	}

	mockProductRepo.On("GetProductByID", ctx, int32(50)).Return(entity.Product{}, pgx.ErrNoRows).Once()

	_, err := pu.Checkout(ctx, req)

	mockProductRepo.AssertExpectations(t)
	mockPromoRepo.AssertNotCalled(t, "GetProductPromotionsByProductID", mock.Anything)
	mockProductRepo.AssertNotCalled(t, "UpdateProductQuantityByID", mock.Anything, mock.Anything)

	if err == nil {
		t.Errorf("Expected error, but got nil")
	}

	httpError, ok := err.(*echo.HTTPError)
	if !ok {
		t.Errorf("Expected echo.HTTPError, but got %T", err)
	}

	if httpError.Code != 404 {
		t.Errorf("Expected status code 404, but got %d", httpError.Code)
	}
}

func TestProductUsecase_Checkout_ProductOutOfStock(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockPromoRepo := new(MockPromotionRepository)

	pu := &ProductUsecase{
		productRepo: mockProductRepo,
		promoRepo:   mockPromoRepo,
	}

	ctx := context.Background()

	req := dto.CheckoutRequest{
		Items: []dto.CheckoutItemData{
			{
				ID:       1,
				Quantity: 20,
			},
		},
	}

	product := entity.Product{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Sku:       "120P90",
		Name:      "Google Home",
		Price:     49.99,
		Quantity:  10,
	}

	mockProductRepo.On("GetProductByID", ctx, int32(1)).Return(product, nil).Once()
	mockPromoRepo.On("GetProductPromotionsByProductID", ctx, int32(1)).Return([]entity.ProductPromotion{}, nil).Once()

	_, err := pu.Checkout(ctx, req)

	mockProductRepo.AssertExpectations(t)
	mockPromoRepo.AssertNotCalled(t, "GetProductPromotionsByProductID", mock.Anything)
	mockProductRepo.AssertNotCalled(t, "UpdateProductQuantityByID", mock.Anything, mock.Anything)

	if err == nil {
		t.Errorf("Expected error, but got nil")
	}

	httpError, ok := err.(*echo.HTTPError)
	if !ok {
		t.Errorf("Expected echo.HTTPError, but got %T", err)
	}

	if httpError.Code != 400 {
		t.Errorf("Expected status code 400, but got %d", httpError.Code)
	}
}

func TestProductUsecase_Checkout_InternalServerError(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockPromoRepo := new(MockPromotionRepository)

	pu := &ProductUsecase{
		productRepo: mockProductRepo,
		promoRepo:   mockPromoRepo,
	}

	ctx := context.Background()

	req := dto.CheckoutRequest{
		Items: []dto.CheckoutItemData{
			{
				ID:       1,
				Quantity: 2,
			},
		},
	}

	product := entity.Product{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Sku:       "120P90",
		Name:      "Google Home",
		Price:     49.99,
		Quantity:  10,
	}

	mockProductRepo.On("GetProductByID", ctx, int32(1)).Return(product, errors.New("internal server error")).Once()
	mockPromoRepo.On("GetProductPromotionsByProductID", ctx, int32(1)).Return([]entity.ProductPromotion{}, nil).Once()

	_, err := pu.Checkout(ctx, req)

	mockProductRepo.AssertExpectations(t)
	mockPromoRepo.AssertNotCalled(t, "GetProductPromotionsByProductID", mock.Anything)
	mockProductRepo.AssertNotCalled(t, "UpdateProductQuantityByID", mock.Anything, mock.Anything)

	if err == nil {
		t.Errorf("Expected error, but got nil")
	}

	httpError, ok := err.(*echo.HTTPError)
	if !ok {
		t.Errorf("Expected echo.HTTPError, but got %T", err)
	}

	if httpError.Code != 500 {
		t.Errorf("Expected status code 500, but got %d", httpError.Code)
	}
}

func TestProductUsecase_Checkout_PromotionInternalServerError(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockPromoRepo := new(MockPromotionRepository)

	pu := &ProductUsecase{
		productRepo: mockProductRepo,
		promoRepo:   mockPromoRepo,
	}

	ctx := context.Background()

	req := dto.CheckoutRequest{
		Items: []dto.CheckoutItemData{
			{
				ID:       1,
				Quantity: 2,
			},
		},
	}

	product := entity.Product{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Sku:       "120P90",
		Name:      "Google Home",
		Price:     49.99,
		Quantity:  10,
	}

	mockProductRepo.On("GetProductByID", ctx, int32(1)).Return(product, nil).Once()
	mockPromoRepo.On("GetProductPromotionsByProductID", ctx, int32(1)).Return([]entity.ProductPromotion{}, errors.New("internal server error")).Once()

	_, err := pu.Checkout(ctx, req)

	mockProductRepo.AssertExpectations(t)
	mockPromoRepo.AssertExpectations(t)
	mockProductRepo.AssertNotCalled(t, "UpdateProductQuantityByID", mock.Anything, mock.Anything)

	if err == nil {
		t.Errorf("Expected error, but got nil")
	}

	httpError, ok := err.(*echo.HTTPError)
	if !ok {
		t.Errorf("Expected echo.HTTPError, but got %T", err)
	}

	if httpError.Code != 500 {
		t.Errorf("Expected status code 500, but got %d", httpError.Code)
	}
}

func TestProductUsecase_ApplyPromotionRules_FREE_ITEM(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockPromoRepo := new(MockPromotionRepository)

	pu := &ProductUsecase{
		productRepo: mockProductRepo,
		promoRepo:   mockPromoRepo,
	}

	ctx := context.Background()

	product := entity.Product{
		ID:        2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Sku:       "43N23P",
		Name:      "Mackbook Pro",
		Price:     5399.99,
		Quantity:  1,
	}

	freeProductID := int32(4)
	freeQuantity := int32(1)

	promotions := []entity.ProductPromotion{
		{
			PromotionType: "FREE_ITEM",
			ID:            1,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			PromotionID:   1,
			ProductID:     2,
			MinQuantity:   1,
			FreeProductID: &freeProductID,
			FreeQuantity:  &freeQuantity,
		},
	}

	freeItem := entity.Product{
		ID:        4,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Sku:       "234234",
		Name:      "Raspberry Pi B",
		Price:     30.00,
		Quantity:  2,
	}

	mockProductRepo.On("GetProductByID", ctx, int32(4)).Return(freeItem, nil).Once()

	checkoutItems, err := pu.ApplyPromotionRules(ctx, product, 1, promotions)

	mockProductRepo.AssertExpectations(t)
	mockPromoRepo.AssertNotCalled(t, "GetProductPromotionsByProductID", mock.Anything)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedCheckoutItems := []entity.CheckoutItem{
		{
			ProductID:   2,
			ProductName: "Mackbook Pro",
			Quantity:    1,
			TotalPrice:  5399.99,
		},
		{
			ProductID:   4,
			ProductName: "Raspberry Pi B",
			Quantity:    1,
			TotalPrice:  0,
		},
	}

	if !reflect.DeepEqual(checkoutItems, expectedCheckoutItems) {
		t.Errorf("Checkout items are not as expected: \nExpected: %v\nGot: %v", expectedCheckoutItems, checkoutItems)
	}
}

func TestProductUsecase_ApplyPromotionRules_BUY_X_PAY_Y(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockPromoRepo := new(MockPromotionRepository)

	pu := &ProductUsecase{
		productRepo: mockProductRepo,
		promoRepo:   mockPromoRepo,
	}

	ctx := context.Background()

	product := entity.Product{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Sku:       "120P90",
		Name:      "Google Home",
		Price:     49.99,
		Quantity:  30,
	}

	payY := int32(2)

	promotions := []entity.ProductPromotion{
		{
			PromotionType: "BUY_X_PAY_Y",
			ID:            2,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			PromotionID:   2,
			ProductID:     1,
			MinQuantity:   3,
			PayY:          &payY,
		},
	}

	checkoutItems, err := pu.ApplyPromotionRules(ctx, product, 3, promotions)

	mockProductRepo.AssertNotCalled(t, "GetProductByID", mock.Anything)
	mockPromoRepo.AssertNotCalled(t, "GetProductPromotionsByProductID", mock.Anything)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedCheckoutItems := []entity.CheckoutItem{
		{
			ProductID:   1,
			ProductName: "Google Home",
			Quantity:    3,
			TotalPrice:  99.98,
		},
	}

	if !reflect.DeepEqual(checkoutItems, expectedCheckoutItems) {
		t.Errorf("Checkout items are not as expected: \nExpected: %v\nGot: %v", expectedCheckoutItems, checkoutItems)
	}
}

func TestProductUsecase_ApplyPromotionRules_BULK_DISCOUNT(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockPromoRepo := new(MockPromotionRepository)

	pu := &ProductUsecase{
		productRepo: mockProductRepo,
		promoRepo:   mockPromoRepo,
	}

	ctx := context.Background()

	product := entity.Product{
		ID:        3,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Sku:       "A304SD",
		Name:      "Alexa Speaker",
		Price:     109.50,
		Quantity:  10,
	}

	discount := float64(0.1)

	promotions := []entity.ProductPromotion{
		{
			PromotionType: "BULK_DISCOUNT",
			ID:            3,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			PromotionID:   3,
			ProductID:     3,
			MinQuantity:   3,
			Discount:      &discount,
		},
	}

	checkoutItems, err := pu.ApplyPromotionRules(ctx, product, 3, promotions)

	mockProductRepo.AssertNotCalled(t, "GetProductByID", mock.Anything)
	mockPromoRepo.AssertNotCalled(t, "GetProductPromotionsByProductID", mock.Anything)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedCheckoutItems := []entity.CheckoutItem{
		{
			ProductID:   3,
			ProductName: "Alexa Speaker",
			Quantity:    3,
			TotalPrice:  295.65,
		},
	}

	if !reflect.DeepEqual(checkoutItems, expectedCheckoutItems) {
		t.Errorf("Checkout items are not as expected: \nExpected: %v\nGot: %v", expectedCheckoutItems, checkoutItems)
	}
}

func TestProductUsecase_ApplyPromotionRules_Default(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockPromoRepo := new(MockPromotionRepository)

	pu := &ProductUsecase{
		productRepo: mockProductRepo,
		promoRepo:   mockPromoRepo,
	}

	ctx := context.Background()

	product := entity.Product{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Sku:       "120P90",
		Name:      "Google Home",
		Price:     49.99,
		Quantity:  10,
	}

	promotions := []entity.ProductPromotion{
		{
			PromotionType: "UNKNOWN",
			ID:            1,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			PromotionID:   1,
			ProductID:     1,
			MinQuantity:   2,
		},
	}

	checkoutItems, err := pu.ApplyPromotionRules(ctx, product, 2, promotions)

	mockProductRepo.AssertNotCalled(t, "GetProductByID", mock.Anything)
	mockPromoRepo.AssertNotCalled(t, "GetProductPromotionsByProductID", mock.Anything)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedCheckoutItems := []entity.CheckoutItem{
		{
			ProductID:   1,
			ProductName: "Google Home",
			Quantity:    2,
			TotalPrice:  99.98,
		},
	}

	if !reflect.DeepEqual(checkoutItems, expectedCheckoutItems) {
		t.Errorf("Checkout items are not as expected: \nExpected: %v\nGot: %v", expectedCheckoutItems, checkoutItems)
	}
}

func TestProductUsecase_ApplyPromotionRules_FreeItemNotFound(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockPromoRepo := new(MockPromotionRepository)

	pu := &ProductUsecase{
		productRepo: mockProductRepo,
		promoRepo:   mockPromoRepo,
	}

	ctx := context.Background()

	product := entity.Product{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Sku:       "43N23P",
		Name:      "Mackbook Pro",
		Price:     5399.99,
		Quantity:  5,
	}

	freeProductID := int32(4)
	freeQuantity := int32(1)

	promotions := []entity.ProductPromotion{
		{
			PromotionType: "FREE_ITEM",
			ID:            1,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			PromotionID:   1,
			ProductID:     1,
			MinQuantity:   1,
			FreeProductID: &freeProductID,
			FreeQuantity:  &freeQuantity,
		},
	}

	mockProductRepo.On("GetProductByID", ctx, int32(4)).Return(entity.Product{}, pgx.ErrNoRows).Once()

	checkoutItems, err := pu.ApplyPromotionRules(ctx, product, 1, promotions)

	mockProductRepo.AssertExpectations(t)
	mockPromoRepo.AssertNotCalled(t, "GetProductPromotionsByProductID", mock.Anything)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedCheckoutItems := []entity.CheckoutItem{
		{
			ProductID:   1,
			ProductName: "Mackbook Pro",
			Quantity:    1,
			TotalPrice:  5399.99,
		},
	}

	if !reflect.DeepEqual(checkoutItems, expectedCheckoutItems) {
		t.Errorf("Checkout items are not as expected: \nExpected: %v\nGot: %v", expectedCheckoutItems, checkoutItems)
	}
}

func TestProductUsecase_ApplyPromotionRules_FreeItemInternalServerError(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockPromoRepo := new(MockPromotionRepository)

	pu := &ProductUsecase{
		productRepo: mockProductRepo,
		promoRepo:   mockPromoRepo,
	}

	ctx := context.Background()

	product := entity.Product{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Sku:       "120P90",
		Name:      "Google Home",
		Price:     49.99,
		Quantity:  10,
	}

	freeProductID := int32(2)
	freeQuantity := int32(1)

	promotions := []entity.ProductPromotion{
		{
			PromotionType: "FREE_ITEM",
			ID:            1,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			PromotionID:   1,
			ProductID:     1,
			MinQuantity:   2,
			FreeProductID: &freeProductID,
			FreeQuantity:  &freeQuantity,
		},
	}

	mockProductRepo.On("GetProductByID", ctx, int32(2)).Return(entity.Product{}, errors.New("internal server error")).Once()

	_, err := pu.ApplyPromotionRules(ctx, product, 2, promotions)

	mockProductRepo.AssertExpectations(t)
	mockPromoRepo.AssertNotCalled(t, "GetProductPromotionsByProductID", mock.Anything)

	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
	assert.Error(t, err)
}
