package models

import (
	"time"

	"gorm.io/gorm"
)

type InvoiceStatus string

const (
	InvoicePending   InvoiceStatus = "pending"
	InvoiceDelivered InvoiceStatus = "delivered"
	InvoiceCancelled InvoiceStatus = "cancelled"
)

type Invoice struct {
	ID        int            `gorm:"primaryKey;column:id" json:"id"`
	UserID    *int           `gorm:"column:user_id" json:"user_id"`
	ProductID *int           `gorm:"column:product_id" json:"product_id"`
	Status    InvoiceStatus  `gorm:"column:status;type:invoice_status;default:pending" json:"status"`
	Note      *string        `gorm:"column:note" json:"note"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`

	User    *User    `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"-"`
	Product *Product `gorm:"foreignKey:ProductID;constraint:OnDelete:SET NULL" json:"-"`
}

func (Invoice) TableName() string {
	return "invoices"
}