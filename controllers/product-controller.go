package controllers

import (
	"Savannah_Screening_Test/config"
	"Savannah_Screening_Test/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	CategoryID  string  `json:"category_id"` //could be nullable in the DB but let me force for it's provision
}

func CreateProduct(c *gin.Context) {
	var request CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := entity.Product{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
	}

	if request.CategoryID != "" {
		categoryUUID, err := ParseUUID(request.CategoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CategoryID is invalid"})
			return
		}
		product.CategoryID = &categoryUUID
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"product": product})
}

func GetProducts(c *gin.Context) {
	categoryID := c.Query("category_id")
	var products []entity.Product

	query := config.DB.Preload("Category")
	if categoryID != "" {
		valCategoryUUID, err := ParseUUID(categoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CategoryID is invalid"})
			return
		}
		query = query.Where("category_id = ?", valCategoryUUID)
	}

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	c.JSON(http.StatusOK, products)
}
