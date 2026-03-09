package database

import "time"

type UserRole string

const (
	UserRoleRoot    UserRole = "root"
	UserRoleAdmin   UserRole = "admin"
	UserRoleManager UserRole = "manager"
	UserRoleMember  UserRole = "member"
)

type User struct {
	Model
	Name        string     `db:"name" json:"name"`
	Email       string     `db:"email" json:"email"`
	Phone       string     `db:"phone" json:"phone"`
	Role        UserRole   `db:"role" json:"role"`
	City        string     `db:"city" json:"city"`
	Governorate string     `db:"governorate" json:"governorate"`
	Address     string     `db:"address" json:"address"`
	Birthdate   *time.Time `db:"birthdate" json:"birthdate"`
}
