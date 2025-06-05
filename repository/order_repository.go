package repository

import (
	"Savannah_Screening_Test/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *entity.Order) error
	CreateOrderItem(item *entity.OrderItem) error
	FindProductByID(id uuid.UUID) (*entity.Product, error)
	GetOrdersByCustomerID(customerID uuid.UUID) ([]entity.Order, error)
	WithTx(tx *gorm.DB) OrderRepository
	BeginTx() *gorm.DB
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) BeginTx() *gorm.DB {
	return r.db.Begin()
}

func (r *orderRepository) WithTx(tx *gorm.DB) OrderRepository {
	return &orderRepository{db: tx}
}

func (r *orderRepository) CreateOrder(order *entity.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) CreateOrderItem(item *entity.OrderItem) error {
	return r.db.Create(item).Error
}

func (r *orderRepository) FindProductByID(id uuid.UUID) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *orderRepository) GetOrdersByCustomerID(customerID uuid.UUID) ([]entity.Order, error) {
	var orders []entity.Order
	err := r.db.Preload("Items.Product").Where("customer_id = ?", customerID).Find(&orders).Error
	return orders, err
}
