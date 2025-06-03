package controllers

import (
	"Savannah_Screening_Test/clients"
	"Savannah_Screening_Test/config"
	"Savannah_Screening_Test/entity"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
)

type OrderItemRequest struct {
	//OrderID   string `json:"order_id" binding:"required"`
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}

type CreateOrderRequest struct {
	//CustomerID string             `json:"customer_id" binding:"required"`
	OrderItems []OrderItemRequest `json:"order_items" binding:"required,dive"` //must not be nil, also validate the items of the list{slice}
}

func CreateOrder(c *gin.Context) {
	var request CreateOrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := c.MustGet("user").(jwt.MapClaims)
	customerID := claims["sub"].(string)
	customerName := claims["name"].(string)
	fmt.Println("customer making order is:: " + customerID)

	customerUUID, err := ParseUUID(customerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"}) //validate customer ID itself in DB
		return
	}

	order := entity.Order{
		CustomerID: customerUUID,
	}

	tx := config.DB.Begin()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order :: FK -> customerID"})
		return
	}

	var orderTotal float64
	for _, orderItem := range request.OrderItems {
		productUUID, err := ParseUUID(orderItem.ProductID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID ->Not a valid UUID"})
			return
		}

		var product entity.Product
		if err := tx.First(&product, "id = ?", productUUID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found :: FK -> productID"})
			return
		}

		orderTotal = orderTotal + (product.Price * float64(orderItem.Quantity))

		orderItem := entity.OrderItem{
			OrderID:   order.ID,
			ProductID: product.ID,
			Quantity:  orderItem.Quantity,
			Price:     product.Price,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create orderItem"})
			return
		}
	}
	tx.Commit()

	customerPhoneNumber := os.Getenv("AFRICASTALKING_SANDBOX_CLIENT_NUMBER")
	customerOrderSummary := fmt.Sprintf("Order with ID %s created for customer %s", order.ID, customerName)
	orderSummary := fmt.Sprintf("New order by %s (customer ID: %s):\nTotal: %f", customerName, customerID, orderTotal)

	clients.SendSMSAsync(customerPhoneNumber, customerOrderSummary)
	clients.SendAdminEmailAsync("New Order Received", orderSummary)

	c.JSON(http.StatusCreated, gin.H{"message": "Order created", "order_id": order.ID})
}

func GetOrdersByCustomer(c *gin.Context) {
	claims := c.MustGet("user").(jwt.MapClaims)
	customerID := claims["sub"].(string)

	customerUUID, err := ParseUUID(customerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var orders []entity.Order
	if err := config.DB.Preload("Items.Product").Where("customer_id = ?", customerUUID).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}
