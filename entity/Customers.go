package entity

import (
	"github.com/google/uuid"
	"time"
)

type Customer struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Phone     string    `gorm:"type:varchar(15);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Orders    []Order   `gorm:"foreignKey:CustomerID"`
}
