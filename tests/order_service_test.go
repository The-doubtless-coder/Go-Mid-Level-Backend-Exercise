package tests

import (
	"Savannah_Screening_Test/dtos"
	"Savannah_Screening_Test/entity"
	"Savannah_Screening_Test/service"
	"Savannah_Screening_Test/tests/mocks"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestCreateOrder_Success(t *testing.T) {
//	mockRepo := new(mocks.MockOrderRepository)
//	mockTx := &gorm.DB{}
//	mockRepo.Tx = mockTx
//
//	dummyTx := &gorm.DB{} // just a placeholder to avoid nil panic
//	mockRepo.On("BeginTx").Return(dummyTx)
//	mockRepo.On("WithTx", dummyTx).Return(mockRepo)
//
//	//mockRepo.On("BeginTx").Return(mockTx)
//	//mockRepo.On("WithTx", mockTx).Return(mockRepo)
//
//	customerID := uuid.New()
//	productID := uuid.New()
//	orderID := uuid.New()
//
//	mockRepo.On("CreateOrder", mock.AnythingOfType("*entity.Order")).Run(func(args mock.Arguments) {
//		arg := args.Get(0).(*entity.Order)
//		arg.ID = orderID
//	}).Return(nil)
//
//	mockRepo.On("FindProductByID", productID).Return(&entity.Product{
//		ID:    productID,
//		Price: 20.0,
//	}, nil)
//
//	mockRepo.On("CreateOrderItem", mock.AnythingOfType("*entity.OrderItem")).Return(nil)
//
//	svc := service.NewOrderService(mockRepo)
//
//	req := dtos.CreateOrderRequest{
//		OrderItems: []dtos.OrderItemRequest{
//			{
//				ProductID: productID.String(),
//				Quantity:  2,
//			},
//		},
//	}
//
//	returnedID, err := svc.CreateOrder(customerID, "Test User", req)
//
//	assert.NoError(t, err)
//	assert.Equal(t, orderID, returnedID)
//	mockRepo.AssertExpectations(t)
//}

func TestGetOrdersByCustomerID(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	customerID := uuid.New()
	mockOrders := []entity.Order{
		{ID: uuid.New(), CustomerID: customerID},
	}

	mockRepo.On("GetOrdersByCustomerID", customerID).Return(mockOrders, nil)

	svc := service.NewOrderService(mockRepo)

	orders, err := svc.GetOrdersByCustomerID(customerID)

	assert.NoError(t, err)
	assert.Equal(t, mockOrders, orders)
	mockRepo.AssertExpectations(t)
}

func TestCreateOrder_InvalidUUID(t *testing.T) {
	repo := new(mocks.MockOrderRepository)
	svc := service.NewOrderService(repo)

	req := dtos.CreateOrderRequest{
		OrderItems: []dtos.OrderItemRequest{{ProductID: "invalid-uuid", Quantity: 1}},
	}

	_, err := svc.CreateOrder(uuid.New(), "Test User", req)
	fmt.Println(err)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid product ID")
}

//func TestCreateOrder_ProductNotFound(t *testing.T) {
//	repo := new(mocks.MockOrderRepository)
//	repo.On("BeginTx").Return(&gorm.DB{})
//	repo.On("WithTx", mock.Anything).Return(repo)
//	repo.On("CreateOrder", mock.Anything).Return(nil)
//	repo.On("FindProductByID", mock.Anything).Return(nil, gorm.ErrRecordNotFound)
//
//	req := dtos.CreateOrderRequest{
//		OrderItems: []dtos.OrderItemRequest{{ProductID: uuid.New().String(), Quantity: 1}},
//	}
//
//	svc := service.NewOrderService(repo)
//	_, err := svc.CreateOrder(uuid.New(), "Test User", req)
//	assert.Error(t, err)
//	assert.Contains(t, err.Error(), "product not found")
//}
