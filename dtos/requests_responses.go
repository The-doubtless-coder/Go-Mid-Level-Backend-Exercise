package dtos

import (
	"github.com/google/uuid"
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

// CreateProductRequest represents product creation payload
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	CategoryID  string  `json:"category_id"`
}

type OrderItemRequest struct {
	//OrderID   string `json:"order_id" binding:"required"`
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}

type CreateOrderRequest struct {
	//CustomerID string             `json:"customer_id" binding:"required"`
	OrderItems []OrderItemRequest `json:"order_items" binding:"required,dive"` //must not be nil, also validate the items of the list{slice}
}

type OrderItemsResponse struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type OrderResponse struct {
	OrderID    string               `json:"order_id"`
	CustomerID string               `json:"customer_id"`
	OrderDate  string               `json:"order_date"`
	OrderItems []OrderItemsResponse `json:"order_items"`
}

// CreateProductResponse represents product response data
type CreateProductResponse struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  *uuid.UUID
	CreatedAt   string `json:"created_at"`
}
type CreateCategoryResponse struct {
	CategoryID       string     `json:"category_id"`
	CategoryName     string     `json:"category_name"`
	CategoryParentID *uuid.UUID `json:"category_parent_id"`
	CreatedAt        string     `json:"created_at"`
}

type GetAllProductsResponse struct {
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	Price        float64 `json:"price"`
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type EmailRequest struct {
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOi..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOi..."`
	ExpiresIn    int    `json:"expires_in" example:"3600"`
}
