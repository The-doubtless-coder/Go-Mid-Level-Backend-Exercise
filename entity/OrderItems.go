package entity

import "github.com/google/uuid"

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null"`
	Order     Order     `gorm:"foreignKey:OrderID"`
	ProductID uuid.UUID `gorm:"type:uuid"`
	Product   Product   `gorm:"foreignKey:ProductID"`
	Quantity  int       `gorm:"not null"`
	Price     float64   `gorm:"type:numeric(10,2);not null"`
}
