package mocks

import (
	"Savannah_Screening_Test/dtos"
	"Savannah_Screening_Test/entity"
	"github.com/google/uuid"
)

// --- Category Mocks ---
type mockCategoryRepo struct {
	CreateFunc   func(*entity.Category) error
	FindByIDFunc func(uuid.UUID) (*entity.Category, error)
}

func (m *mockCategoryRepo) Create(cat *entity.Category) error {
	return m.CreateFunc(cat)
}
func (m *mockCategoryRepo) FindByID(id uuid.UUID) (*entity.Category, error) {
	return m.FindByIDFunc(id)
}

// --- Product Mocks ---
type mockProductRepo struct {
	CreateFunc                     func(*entity.Product) error
	FindAllFunc                    func(*uuid.UUID) ([]entity.Product, error)
	GetAveragePricePerCategoryFunc func() ([]dtos.AvgPriceResponse, error)
}

func (m *mockProductRepo) Create(p *entity.Product) error {
	return m.CreateFunc(p)
}
func (m *mockProductRepo) FindAll(categoryID *uuid.UUID) ([]entity.Product, error) {
	return m.FindAllFunc(categoryID)
}
func (m *mockProductRepo) GetAveragePricePerCategory() ([]dtos.AvgPriceResponse, error) {
	return m.GetAveragePricePerCategoryFunc()
}
