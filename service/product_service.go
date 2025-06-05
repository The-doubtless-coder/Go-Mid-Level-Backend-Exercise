package service

import (
	"Savannah_Screening_Test/dtos"
	"Savannah_Screening_Test/entity"
	"Savannah_Screening_Test/repository"
	"github.com/google/uuid"
)

type ProductService interface {
	Create(product *entity.Product) error
	GetAll(categoryID *uuid.UUID) ([]entity.Product, error)
	GetAveragePricePerCategory() ([]dtos.AvgPriceResponse, error)
}

type productService struct {
	repo repository.ProductRepository
}

//type ProductService interface {
//	GetAveragePricePerCategory() ([]dtos.AvgPriceResponse, error)
//}
//type productService struct {
//	repo repository.ProductRepository
//}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetAveragePricePerCategory() ([]dtos.AvgPriceResponse, error) {
	return s.repo.GetAveragePricePerCategory()
}

func (s *productService) Create(product *entity.Product) error {
	return s.repo.Create(product)
}

func (s *productService) GetAll(categoryID *uuid.UUID) ([]entity.Product, error) {
	return s.repo.FindAll(categoryID)
}
