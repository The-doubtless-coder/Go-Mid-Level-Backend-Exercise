package service

import (
	"Savannah_Screening_Test/clients"
	"Savannah_Screening_Test/dtos"
	"Savannah_Screening_Test/entity"
	"Savannah_Screening_Test/repository"
	"fmt"
	"github.com/google/uuid"
	"os"
)

type OrderService interface {
	CreateOrder(customerID uuid.UUID, customerName string, req dtos.CreateOrderRequest) (uuid.UUID, error)
	GetOrdersByCustomerID(customerID uuid.UUID) ([]entity.Order, error)
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) CreateOrder(customerID uuid.UUID, customerName string, req dtos.CreateOrderRequest) (uuid.UUID, error) {
	productIDs := make([]uuid.UUID, 0, len(req.OrderItems))
	for _, item := range req.OrderItems {
		productID, err := uuid.Parse(item.ProductID)
		if err != nil {
			return uuid.Nil, fmt.Errorf("invalid product ID: %w", err)
		}
		productIDs = append(productIDs, productID)
	} //todo: makes sure with invalid UUIDs passes

	order := entity.Order{CustomerID: customerID}
	tx := s.repo.BeginTx()

	if tx == nil {
		return uuid.Nil, fmt.Errorf("failed to begin transaction")
	}

	r := s.repo.WithTx(tx)

	if err := r.CreateOrder(&order); err != nil {
		if tx != nil {
			tx.Rollback()
		}
		return uuid.Nil, fmt.Errorf("failed to create order: %w", err)
	}

	var total float64
	for _, item := range req.OrderItems {
		productID, err := uuid.Parse(item.ProductID) //todo: no need for this
		if err != nil {
			if tx != nil {
				tx.Rollback()
			}
			return uuid.Nil, fmt.Errorf("invalid product ID: %w", err)
		}

		product, err := r.FindProductByID(productID) //todo: use my slice of IDS now
		if err != nil {
			if tx != nil {
				tx.Rollback()
			}
			return uuid.Nil, fmt.Errorf("product not found: %w", err)
		}

		total += product.Price * float64(item.Quantity)
		orderItem := entity.OrderItem{
			OrderID:   order.ID,
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}
		if err := r.CreateOrderItem(&orderItem); err != nil {
			if tx != nil {
				tx.Rollback()
			}
			return uuid.Nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}
	if err := tx.Commit().Error; err != nil {
		return uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	tx.Commit()

	sendNotification(customerName, order.ID, customerID.String(), total)
	return order.ID, nil
}

func (s *orderService) GetOrdersByCustomerID(customerID uuid.UUID) ([]entity.Order, error) {
	return s.repo.GetOrdersByCustomerID(customerID)
}

func sendNotification(name string, orderID uuid.UUID, customerID string, total float64) {
	msg := fmt.Sprintf("New order by %s (ID: %s):\nTotal: %.2f", name, customerID, total)
	clients.SendSMSAsync(os.Getenv("AFRICASTALKING_SANDBOX_CLIENT_NUMBER"), fmt.Sprintf("Order %s created", orderID))
	clients.SendAdminEmailAsync("New Order Received", msg)
}
