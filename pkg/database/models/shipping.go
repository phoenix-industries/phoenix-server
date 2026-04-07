package models

import (
	"time"

	"gorm.io/gorm"
)

type Shipping struct {
	ID        int            `gorm:"primaryKey;column:id" json:"id"`
	FromID    *int           `gorm:"column:from_id" json:"from_id"`
	ToID      *int           `gorm:"column:to_id" json:"to_id"`
	Fee       float64        `gorm:"column:fee;type:numeric(10,2);not null" json:"fee"`
	Location  string         `gorm:"column:location;not null" json:"location"`
	Note      *string        `gorm:"column:note" json:"note"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`

	FromUser *User `gorm:"foreignKey:FromID;constraint:OnDelete:SET NULL" json:"-"`
	ToUser   *User `gorm:"foreignKey:ToID;constraint:OnDelete:SET NULL" json:"-"`
}

func (Shipping) TableName() string {
	return "shippings"
}