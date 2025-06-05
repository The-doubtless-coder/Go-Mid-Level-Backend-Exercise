package handler

import (
	"Savannah_Screening_Test/dtos"
	"Savannah_Screening_Test/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (oc *OrderHandler) CreateOrder(c *gin.Context) {
	var request dtos.CreateOrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := c.MustGet("user").(jwt.MapClaims)
	customerID := claims["sub"].(string)
	customerName := claims["name"].(string)

	customerUUID, err := ParseUUID(customerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderID, err := oc.orderService.CreateOrder(customerUUID, customerName, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Order created", "order_id": orderID})
}

func (oc *OrderHandler) GetOrdersByCustomer(c *gin.Context) {
	claims := c.MustGet("user").(jwt.MapClaims)
	customerID := claims["sub"].(string)

	customerUUID, err := ParseUUID(customerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orders, err := oc.orderService.GetOrdersByCustomerID(customerUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, MapOrdersToResponses(orders))
}
