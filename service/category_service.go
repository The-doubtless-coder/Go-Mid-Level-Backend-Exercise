package service

import (
	"Savannah_Screening_Test/dtos"
	"Savannah_Screening_Test/entity"
	"Savannah_Screening_Test/repository"
	"fmt"
	"github.com/google/uuid"
)

type CategoryService interface {
	CreateCategory(req dtos.CreateCategoryRequest) (*entity.Category, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo}
}

func (s *categoryService) CreateCategory(req dtos.CreateCategoryRequest) (*entity.Category, error) {
	category := &entity.Category{
		Name: req.Name,
	}

	if req.ParentID != "" {
		parentUUID, err := uuid.Parse(req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent UUID")
		}
		_, err = s.repo.FindByID(parentUUID)
		if err != nil {
			return nil, fmt.Errorf("parent category not found")
		}
		category.ParentID = &parentUUID
	}

	if err := s.repo.Create(category); err != nil {
		return nil, fmt.Errorf("failed to create category")
	}
	return category, nil
}
