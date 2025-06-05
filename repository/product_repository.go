package repository

import (
	"Savannah_Screening_Test/dtos"
	"Savannah_Screening_Test/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *entity.Product) error
	FindAll(categoryID *uuid.UUID) ([]entity.Product, error)
	GetAveragePricePerCategory() ([]dtos.AvgPriceResponse, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) Create(product *entity.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindAll(categoryID *uuid.UUID) ([]entity.Product, error) {
	var products []entity.Product
	query := r.db.Preload("Category")
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

//type ProductRepository interface {
//	GetAveragePricePerCategory() ([]dtos.AvgPriceResponse, error)
//}

//type productRepository struct {
//	db *gorm.DB
//}

//func NewProductRepository(db *gorm.DB) ProductRepository {
//	return &productRepository{db: db}
//}

func (r *productRepository) GetAveragePricePerCategory() ([]dtos.AvgPriceResponse, error) {
	var results []dtos.AvgPriceResponse
	query := `SELECT c.id AS category_id, c.name AS category_name, COALESCE(AVG(p.price), 0) AS average_price
              FROM categories c
              LEFT JOIN products p ON p.category_id = c.id
              GROUP BY c.id, c.name
              ORDER BY c.name`

	if err := r.db.Raw(query).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}
