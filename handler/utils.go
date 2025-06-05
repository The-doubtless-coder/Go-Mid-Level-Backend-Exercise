package handler

import (
	"Savannah_Screening_Test/dtos"
	"Savannah_Screening_Test/entity"
	"github.com/google/uuid"
	"time"
)

func ParseUUID(input string) (uuid.UUID, error) {
	return uuid.Parse(input)
}

func MapOrdersToResponses(orders []entity.Order) []dtos.OrderResponse {
	responses := make([]dtos.OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = MapOrderToResponse(order)
	}
	return responses
}

func MapOrderToResponse(order entity.Order) dtos.OrderResponse {
	var ItemsList []dtos.OrderItemsResponse
	for _, item := range order.Items {
		itemsResponse := dtos.OrderItemsResponse{
			ProductID:   item.ProductID.String(),
			ProductName: item.Product.Name,
			Price:       item.Product.Price,
			Quantity:    item.Quantity,
		}
		ItemsList = append(ItemsList, itemsResponse)
	}

	return dtos.OrderResponse{
		OrderID:    order.ID.String(),
		CustomerID: order.CustomerID.String(), // assume FullName() string method on Customer
		OrderDate:  order.OrderDate.String(),
		OrderItems: ItemsList,
	}
}

func MapProductToResponse(product *entity.Product) dtos.CreateProductResponse {
	return dtos.CreateProductResponse{
		ProductID:   product.ID.String(),
		ProductName: product.Name,
		Description: product.Description,
		Price:       product.Price,
		CategoryID:  product.CategoryID,
		CreatedAt:   product.CreatedAt.String(),
	}
}

type Category struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string     `gorm:"type:varchar(100);not null"`
	ParentID  *uuid.UUID `gorm:"type:uuid;default:null"`
	Parent    *Category  `gorm:"foreignKey:ParentID"`
	Children  []Category `gorm:"foreignKey:ParentID"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
}

func MapCategoryToResponse(category *entity.Category) dtos.CreateCategoryResponse {
	return dtos.CreateCategoryResponse{
		CategoryID:       category.ID.String(),
		CategoryName:     category.Name,
		CategoryParentID: category.ParentID,
		CreatedAt:        category.CreatedAt.String(),
	}
}

type GetAllProductsResponse struct {
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	Price        float64 `json:"price"`
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
}

func MapAllProductsToResponses(entity []entity.Product) []dtos.GetAllProductsResponse {
	var productResponses []dtos.GetAllProductsResponse
	for _, product := range entity {
		dto := dtos.GetAllProductsResponse{
			ProductID:    product.ID.String(),
			ProductName:  product.Name,
			Price:        product.Price,
			CategoryID:   product.CategoryID.String(),
			CategoryName: product.Category.Name,
		}
		productResponses = append(productResponses, dto)
	}
	return productResponses
}
