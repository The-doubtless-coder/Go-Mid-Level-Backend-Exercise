package entity

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:varchar(255);not null"`
	Price       float64   `gorm:"type:numeric(10,2);not null"`
	CategoryID  *uuid.UUID
	Category    *Category `gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
