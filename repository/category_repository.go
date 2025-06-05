package repository

import (
	"Savannah_Screening_Test/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *entity.Category) error
	FindByID(id uuid.UUID) (*entity.Category, error)
}

type categoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepo{db}
}

func (r *categoryRepo) Create(category *entity.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepo) FindByID(id uuid.UUID) (*entity.Category, error) {
	var cat entity.Category
	if err := r.db.First(&cat, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}
