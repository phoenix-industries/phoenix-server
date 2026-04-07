package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          int            `gorm:"primaryKey;column:id" json:"id"`
	UserID      *int           `gorm:"column:user_id" json:"user_id"`
	Name        string         `gorm:"column:name;not null" json:"name"`
	CategoryID  *int           `gorm:"column:category" json:"category_id"` // matches schema column name "category"
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`

	// Relations
	User     *User           `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"-"`
	Category *ProductCategory `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL" json:"-"`
	Images   []ProductImage  `gorm:"foreignKey:ProductID;constraint:OnDelete:SET NULL" json:"-"`
	Tags     []ProductTag    `gorm:"foreignKey:ProductID;constraint:OnDelete:SET NULL" json:"-"`
	Reviews  []ProductReview `gorm:"foreignKey:ProductID;constraint:OnDelete:SET NULL" json:"-"`
	Invoices []Invoice       `gorm:"foreignKey:ProductID;constraint:OnDelete:SET NULL" json:"-"`
}

func (Product) TableName() string {
	return "products"
}