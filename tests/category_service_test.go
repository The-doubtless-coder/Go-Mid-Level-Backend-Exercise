package tests

import (
	"Savannah_Screening_Test/dtos"
	"Savannah_Screening_Test/entity"
	"Savannah_Screening_Test/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockRepo struct {
	CreateFunc   func(*entity.Category) error
	FindByIDFunc func(uuid.UUID) (*entity.Category, error)
}

func (m *mockRepo) Create(cat *entity.Category) error {
	return m.CreateFunc(cat)
}
func (m *mockRepo) FindByID(id uuid.UUID) (*entity.Category, error) {
	return m.FindByIDFunc(id)
}

func TestCreateCategory(t *testing.T) {
	mock := &mockRepo{
		CreateFunc: func(cat *entity.Category) error {
			return nil
		},
		FindByIDFunc: func(id uuid.UUID) (*entity.Category, error) {
			return &entity.Category{ID: id, Name: "Parent"}, nil
		},
	}

	svc := service.NewCategoryService(mock)
	req := dtos.CreateCategoryRequest{
		Name:     "Electronics",
		ParentID: uuid.New().String(),
	}

	result, err := svc.CreateCategory(req)
	assert.NoError(t, err)
	assert.Equal(t, "Electronics", result.Name)
}
