// Package models contains database models.
package models

import "time"

type Model struct {
	ID        string     `db:"id" json:"id"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}
