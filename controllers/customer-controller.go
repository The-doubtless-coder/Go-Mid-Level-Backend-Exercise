package controllers

import (
	"Savannah_Screening_Test/config"
	"Savannah_Screening_Test/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateCustomerRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required" validate:"email"`
	Phone string `json:"phone" binding:"required" validate:"phone"`
}

func CreateCustomer(c *gin.Context) {
	var request CreateCustomerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//check if email is unique
	var existing entity.Customer
	if err := config.DB.Where("email = ?", request.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	customer := entity.Customer{
		Name:  request.Name,
		Email: request.Email,
		Phone: request.Phone,
	}

	if err := config.DB.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	//save user in my OIDC provide //pass username (email + password)

	c.JSON(http.StatusCreated, customer)
}
