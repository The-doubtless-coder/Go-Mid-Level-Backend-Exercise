package migrations

import (
	"Savannah_Screening_Test/entity"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.Customer{},
		&entity.Category{},
		&entity.Product{},
		&entity.Order{},
		&entity.OrderItem{},
	)
}
