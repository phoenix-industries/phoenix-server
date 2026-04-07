package models

import (
	"time"

	"gorm.io/gorm"
)

type UserBan struct {
	ID          int            `gorm:"primaryKey;column:id" json:"id"`
	UserID      *int           `gorm:"column:user_id" json:"user_id"`
	ModeratorID *int           `gorm:"column:moderator_id" json:"moderator_id"`
	Reason      string         `gorm:"column:reason;not null" json:"reason"`
	Message     string         `gorm:"column:message;not null" json:"message"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`

	// Relations
	User      *User `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"-"`
	Moderator *User `gorm:"foreignKey:ModeratorID;constraint:OnDelete:SET NULL" json:"-"`
}

func (UserBan) TableName() string {
	return "user_bans"
}