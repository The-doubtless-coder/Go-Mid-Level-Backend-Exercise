package entity

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID         uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerID uuid.UUID   `gorm:"type:uuid;not null"`
	Customer   Customer    `gorm:"foreignKey:CustomerID"`
	OrderDate  time.Time   `gorm:"autoCreateTime"`
	Items      []OrderItem `gorm:"foreignKey:OrderID"`
}
