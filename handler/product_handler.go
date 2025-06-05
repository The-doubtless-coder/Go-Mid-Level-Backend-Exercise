package handler

import (
	"Savannah_Screening_Test/dtos"
	"Savannah_Screening_Test/entity"
	"Savannah_Screening_Test/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(s service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (p *ProductHandler) CreateProduct(c *gin.Context) {
	var request dtos.CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := &entity.Product{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
	}

	if request.CategoryID != "" {
		categoryUUID, err := uuid.Parse(request.CategoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CategoryID is invalid"})
			return
		}
		product.CategoryID = &categoryUUID
	}

	if err := p.service.Create(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"product": MapProductToResponse(product)})
}

func (p *ProductHandler) GetProducts(c *gin.Context) {
	var categoryID *uuid.UUID
	queryParam := c.Query("category_id")
	if queryParam != "" {
		parsed, err := uuid.Parse(queryParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CategoryID is invalid"})
			return
		}
		categoryID = &parsed
	}

	products, err := p.service.GetAll(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	c.JSON(http.StatusOK, MapAllProductsToResponses(products))
}

func (p *ProductHandler) GetAveragePricePerCategoryHandler(c *gin.Context) {
	results, err := p.service.GetAveragePricePerCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch average prices"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results})
}
