package mocks

import (
	"Savannah_Screening_Test/entity"
	"Savannah_Screening_Test/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockOrderRepository struct {
	mock.Mock
	Tx *gorm.DB
}

func (m *MockOrderRepository) CreateOrder(order *entity.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) CreateOrderItem(item *entity.OrderItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockOrderRepository) FindProductByID(id uuid.UUID) (*entity.Product, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *MockOrderRepository) GetOrdersByCustomerID(customerID uuid.UUID) ([]entity.Order, error) {
	args := m.Called(customerID)
	return args.Get(0).([]entity.Order), args.Error(1)
}

func (m *MockOrderRepository) BeginTx() *gorm.DB {
	return m.Tx
}

func (m *MockOrderRepository) WithTx(tx *gorm.DB) repository.OrderRepository {
	args := m.Called(tx)
	return args.Get(0).(repository.OrderRepository)
}
