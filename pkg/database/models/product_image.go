package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductImage struct {
	ID        int            `gorm:"primaryKey;column:id" json:"id"`
	ProductID *int           `gorm:"column:product_id" json:"product_id"`
	Image     string         `gorm:"column:image;not null" json:"image"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`

	Product *Product `gorm:"foreignKey:ProductID;constraint:OnDelete:SET NULL" json:"-"`
}

func (ProductImage) TableName() string {
	return "product_images"
}