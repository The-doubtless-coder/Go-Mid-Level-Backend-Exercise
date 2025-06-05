package tests

import (
	"Savannah_Screening_Test/dtos"
	"Savannah_Screening_Test/entity"
	"Savannah_Screening_Test/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockProductRepo struct {
	CreateFunc  func(*entity.Product) error
	FindAllFunc func(*uuid.UUID) ([]entity.Product, error)
}

func (m *mockProductRepo) Create(product *entity.Product) error {
	return m.CreateFunc(product)
}

func (m *mockProductRepo) FindAll(categoryID *uuid.UUID) ([]entity.Product, error) {
	return m.FindAllFunc(categoryID)
}

func (m *mockProductRepo) GetAveragePricePerCategory() ([]dtos.AvgPriceResponse, error) {
	return m.GetAveragePricePerCategory()
}

func TestCreateProduct(t *testing.T) {
	mockRepo := &mockProductRepo{
		CreateFunc: func(product *entity.Product) error {
			assert.Equal(t, "Test Product", product.Name)
			assert.Equal(t, 100.0, product.Price)
			return nil
		},
	}

	svc := service.NewProductService(mockRepo)

	product := &entity.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.0,
	}

	err := svc.Create(product)
	assert.NoError(t, err)
}

func TestGetAllProducts(t *testing.T) {
	expectedProducts := []entity.Product{
		{
			Name:        "Item A",
			Description: "Desc A",
			Price:       50.0,
		},
		{
			Name:        "Item B",
			Description: "Desc B",
			Price:       75.0,
		},
	}

	mockRepo := &mockProductRepo{
		FindAllFunc: func(categoryID *uuid.UUID) ([]entity.Product, error) {
			assert.Nil(t, categoryID) // simulate no filter
			return expectedProducts, nil
		},
	}

	svc := service.NewProductService(mockRepo)
	products, err := svc.GetAll(nil)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(products))
	assert.Equal(t, expectedProducts[0].Name, products[0].Name)
}

func TestGetAllProductsByCategory(t *testing.T) {
	categoryID := uuid.New()

	mockRepo := &mockProductRepo{
		FindAllFunc: func(catID *uuid.UUID) ([]entity.Product, error) {
			assert.Equal(t, categoryID, *catID)
			return []entity.Product{{Name: "Item 1", Price: 20}}, nil
		},
	}

	svc := service.NewProductService(mockRepo)
	products, err := svc.GetAll(&categoryID)

	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, "Item 1", products[0].Name)
}
