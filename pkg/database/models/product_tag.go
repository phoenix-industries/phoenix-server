package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductTag struct {
	ID        int            `gorm:"primaryKey;column:id" json:"id"`
	ProductID *int           `gorm:"column:product_id" json:"product_id"`
	Tag       string         `gorm:"column:tag;not null" json:"tag"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`

	Product *Product `gorm:"foreignKey:ProductID;constraint:OnDelete:SET NULL" json:"-"`
}

func (ProductTag) TableName() string {
	return "product_tags"
}