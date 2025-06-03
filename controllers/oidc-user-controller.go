package controllers

import (
	"Savannah_Screening_Test/clients"
	"Savannah_Screening_Test/config"
	"Savannah_Screening_Test/dtos"
	"Savannah_Screening_Test/entity"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SignUpRequest struct {
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required" validate:"phone"`
}

func SignUpHandler(c *gin.Context) {
	var request dtos.SignUpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//check if email is unique
	var existing entity.Customer
	if err := config.DB.Where("email = ?", request.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists ->LOCALDB"})
		return
	}

	token, err := clients.GetKeyCloakAdminToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin login failed"})
		return
	}
	fmt.Println("admin token:: " + token)

	customerID, err := clients.CreateUserInKeycloak(request, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	valUUID, err := ParseUUID(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Value returned not a valid UUID"})
	}

	//create other roles depending on your resources
	if err := clients.AssignRoleToUser(customerID, "customer", token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role"})
		return
	}

	fmt.Println("confirm ID in DB " + customerID)
	//use ID to create customer in local table
	customer := entity.Customer{
		ID:    valUUID, //realm ID from keycloak
		Name:  request.Name,
		Email: request.Email,
		Phone: request.Phone,
		//RealmID: customerID,
	}

	// Save to local DB
	if err := config.DB.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created", "customer": customer})
}

//func SignUp(c *gin.Context) {
//	var request dtos.SignUpRequest
//	if err := c.ShouldBindJSON(&request); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Binding failed, check request Body"})
//		return
//	}
//
//	//check if email is unique
//	var existing entity.Customer
//	if err := config.DB.Where("email = ?", request.Email).First(&existing).Error; err == nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists ->LOCALDB"})
//		return
//	}
//
//	token, err := GetKeyCloakAdminToken()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin login failed"})
//		return
//	}
//	fmt.Println("admin token:: " + token)
//
//	customerID, err := createUserInKeycloak(request, token)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	valUUID, err := ParseUUID(customerID)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Value returned not a valid UUID"})
//	}
//
//	fmt.Println("confirm ID in DB " + customerID)
//	//use ID to create customer in local table
//	customer := entity.Customer{
//		ID:    valUUID, //realm ID from keycloak
//		Name:  request.Name,
//		Email: request.Email,
//		Phone: request.Phone,
//		//RealmID: customerID,
//	}
//
//	if err := config.DB.Create(&customer).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
//		return
//	}
//
//	c.JSON(http.StatusCreated, gin.H{"message": "User created", "customer": customer})
//}
