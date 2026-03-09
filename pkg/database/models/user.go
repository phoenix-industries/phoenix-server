package models

import (
	"context"
	"errors"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/phoenix-industries/phoenix-server/pkg/database"
)

type UserRole string

const (
	UserRoleRoot    UserRole = "root"
	UserRoleAdmin   UserRole = "admin"
	UserRoleManager UserRole = "manager"
	UserRoleMember  UserRole = "member"
)

func (r UserRole) String() string {
	return string(r)
}

func (r UserRole) Scan(value any) error {
	switch value.(string) {
	case "root":
		r = UserRoleRoot
	case "admin":
		r = UserRoleAdmin
	case "manager":
		r = UserRoleManager
	case "member":
		r = UserRoleMember
	default:
		return errors.New("invalid user role")
	}
	return nil
}

type User struct {
	Model
	Name        string     `db:"name" json:"name"`
	Email       string     `db:"email" json:"email"`
	Phone       string     `db:"phone" json:"phone"`
	Role        UserRole   `db:"role" json:"role"`
	City        string     `db:"city" json:"city"`
	Governorate string     `db:"governorate" json:"governorate"`
	Address     string     `db:"address" json:"address"`
	Password    string     `db:"password" json:"-"`
	Birthdate   *time.Time `db:"birthdate" json:"birthdate"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.Phone == "" {
		return errors.New("phone is required")
	}
	if u.City == "" {
		return errors.New("city is required")
	}
	if u.Governorate == "" {
		return errors.New("governorate is required")
	}
	if u.Birthdate == nil {
		return errors.New("birthdate is required")
	}
	return nil
}

// InsertUser inserts a new user into the database returning the user id.
func InsertUser(ctx context.Context, db database.DB, user *User) (string, error) {
	if user.Role == "" {
		user.Role = UserRoleMember
	}
	stmt := `
		INSERT INTO users
		(name, email, phone, role, city, governorate, address, password, birthdate)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	if err := db.QueryRow(ctx, stmt, user.Name, user.Email, user.Phone, user.Role, user.City, user.Governorate, user.Address, user.Password, user.Birthdate).Scan(&user.ID); err != nil {
		return "", err
	}
	return user.ID, nil
}

func GetUserByID(ctx context.Context, db database.DB, id string) (*User, error) {
	stmt := `SELECT * FROM users WHERE id = $1`
	var user User
	if err := pgxscan.Get(ctx, db, &user, stmt, id); err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(ctx context.Context, db database.DB, email string) (*User, error) {
	stmt := `SELECT * FROM users WHERE email = $1`
	var user User
	if err := pgxscan.Get(ctx, db, &user, stmt, email); err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByPhone(ctx context.Context, db database.DB, phone string) (*User, error) {
	stmt := `SELECT * FROM users WHERE phone = $1`
	var user User
	if err := pgxscan.Get(ctx, db, &user, stmt, phone); err != nil {
		return nil, err
	}
	return &user, nil
}
