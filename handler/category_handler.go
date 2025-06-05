package handler

import (
	"Savannah_Screening_Test/dtos"
	"Savannah_Screening_Test/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dtos.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//todo: check that the category name is not already there
	//todo: prevents categories with duplicate names

	category, err := h.service.CreateCategory(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, MapCategoryToResponse(category))
}
