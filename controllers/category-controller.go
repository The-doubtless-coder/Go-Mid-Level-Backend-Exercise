package controllers

import (
	"Savannah_Screening_Test/config"
	"Savannah_Screening_Test/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type CreateCategoryRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID string `json:"parent_id"`
}

type AvgPriceResponse struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	AveragePrice float64 `json:"average_price"`
}

func CreateCategory(c *gin.Context) {
	var request CreateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := entity.Category{
		Name: request.Name,
	}

	if request.ParentID != "" {
		parentUUID, err := uuid.Parse(request.ParentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent_id -> Invalid UUID"})
			return
		}

		var parent entity.Category
		if err := config.DB.First(&parent, "id = ?", parentUUID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent category not found"})
			return
		}

		category.ParentID = &parentUUID
	}

	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func GetAveragePricePerCategory(c *gin.Context) {
	var results []AvgPriceResponse
	query := `SELECT c.id AS category_id, c.name AS category_name, COALESCE(AVG(p.price), 0) AS average_price
        		FROM categories c
        		LEFT JOIN products p ON p.category_id = c.id
        		GROUP BY c.id, c.name
        		ORDER BY c.name`

	if err := config.DB.Raw(query).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch average prices"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}
