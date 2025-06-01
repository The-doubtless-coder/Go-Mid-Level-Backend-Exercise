package controllers

import (
	"Savannah_Screening_Test/config"
	"Savannah_Screening_Test/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrderItemRequest struct {
	//OrderID   string `json:"order_id" binding:"required"`
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}

type CreateOrderRequest struct {
	CustomerID string             `json:"customer_id" binding:"required"`
	OrderItems []OrderItemRequest `json:"order_items" binding:"required,dive"` //must not be nil, also validate the items of the list{slice}
}

func CreateOrder(c *gin.Context) {
	var request CreateOrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerUUID, err := ParseUUID(request.CustomerID)
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
	c.JSON(http.StatusCreated, gin.H{"message": "Order created", "order_id": order.ID})
}
