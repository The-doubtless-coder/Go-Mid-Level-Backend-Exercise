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

// CreateOrder godoc
// @Summary Create a new order for the authenticated user
// @Description Create an order for the authenticated customer
// @Tags orders
// @Accept json
// @Produce json
// @Param order body dtos.CreateOrderRequest true "Order Info"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [post]
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

// GetOrdersByCustomer godoc
// @Summary Get all orders for the authenticated customer
// @Description Returns a list of orders made by the logged-in customer
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dtos.OrderResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [get] map[string]string
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
