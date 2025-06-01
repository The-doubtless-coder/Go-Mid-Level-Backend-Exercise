package entity

import (
	"github.com/google/uuid"
	"time"
)

type Category struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string     `gorm:"type:varchar(100);not null"`
	ParentID  *uuid.UUID `gorm:"type:uuid;default:null"`
	Parent    *Category  `gorm:"foreignKey:ParentID"`
	Children  []Category `gorm:"foreignKey:ParentID"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
}
