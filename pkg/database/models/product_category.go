package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductCategory struct {
	ID        int            `gorm:"primaryKey;column:id" json:"id"`
	Category  string         `gorm:"column:category;not null" json:"category"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`

	// Relation
	Products []Product `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL" json:"-"`
}

func (ProductCategory) TableName() string {
	return "product_categories"
}