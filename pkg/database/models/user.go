package models

import (
	"context"
	"errors"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/phoenix-industries/phoenix-server/pkg/auth"
	"github.com/phoenix-industries/phoenix-server/pkg/database"
	"github.com/phoenix-industries/phoenix-server/pkg/validate"
)

type User struct {
	Model
	Name        string    `db:"name" json:"name"`
	Email       string    `db:"email" json:"email"`
	Phone       string    `db:"phone" json:"phone"`
	Role        auth.Role `db:"role" json:"role"`
	Password    string    `db:"password" json:"-"`
	City        string    `db:"city" json:"city"`
	Governorate string    `db:"governorate" json:"governorate"`
	Address     string    `db:"address" json:"address"`
	Birthdate   time.Time `db:"birthdate" json:"birthdate"`
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
	if u.Role == "" {
		return errors.New("role is required")
	}
	if u.City == "" {
		return errors.New("city is required")
	}
	if u.Governorate == "" {
		return errors.New("governorate is required")
	}
	if u.Birthdate.IsZero() {
		return errors.New("birthdate is required")
	}
	if err := validate.Email(u.Email); err != nil {
		return err
	}
	if err := validate.PhoneNumber(u.Phone); err != nil {
		return err
	}
	return nil
}

func UserInsert(ctx context.Context, db database.DB, user *User) error {
	if user.ID == "" {
		return errors.New("id is not set")
	}
	if user.Role == "" {
		user.Role = auth.RoleMember
	}
	if user.Password == "" {
		return errors.New("password is not set")
	}
	stmt := `
		INSERT INTO users
		(id, name, email, phone, role, city, governorate, address, password, birthdate)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	if _, err := db.Exec(ctx, stmt, user.ID, user.Name, user.Email, user.Phone, user.Role, user.City, user.Governorate, user.Address, user.Password, user.Birthdate); err != nil {
		return err
	}
	return nil
}

func UserGetByID(ctx context.Context, db database.DB, id string) (*User, error) {
	stmt := `SELECT * FROM users WHERE id = $1 and deleted_at is null`
	var user User
	if err := pgxscan.Get(ctx, db, &user, stmt, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func UserGetAll(ctx context.Context, db database.DB, limit int, offset int) ([]User, error) {
	stmt := `SELECT * FROM users WHERE deleted_at is null ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	var users []User
	if err := pgxscan.Select(ctx, db, &users, stmt, limit, offset); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return users, nil
}

func UserGetByEmail(ctx context.Context, db database.DB, email string) (*User, error) {
	stmt := `SELECT * FROM users WHERE email = $1 and deleted_at is null`
	var user User
	if err := pgxscan.Get(ctx, db, &user, stmt, email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func UserGetByPhone(ctx context.Context, db database.DB, phone string) (*User, error) {
	stmt := `SELECT * FROM users WHERE phone = $1 and deleted_at is null`
	var user User
	if err := pgxscan.Get(ctx, db, &user, stmt, phone); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func UserExistsWithID(ctx context.Context, db database.DB, id string) (bool, error) {
	stmt := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND deleted_at IS null)`
	var exists bool
	if err := pgxscan.Get(ctx, db, &exists, stmt, id); err != nil {
		return false, err
	}
	return exists, nil
}

func UserExistsWithEmail(ctx context.Context, db database.DB, email string) (bool, error) {
	stmt := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS null)`
	var exists bool
	if err := pgxscan.Get(ctx, db, &exists, stmt, email); err != nil {
		return false, err
	}
	return exists, nil
}

func UserExistsWithPhone(ctx context.Context, db database.DB, phone string) (bool, error) {
	stmt := `SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1 AND deleted_at IS null)`
	var exists bool
	if err := pgxscan.Get(ctx, db, &exists, stmt, phone); err != nil {
		return false, err
	}
	return exists, nil
}

func UserUpdate(ctx context.Context, db database.DB, user *User) error {
	stmt := `
		UPDATE users
		SET name = $2, email = $3, phone = $4, city = $5, governorate = $6, address = $7, birthdate = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	if _, err := db.Exec(ctx, stmt, user.ID, user.Name, user.Email, user.Phone, user.City, user.Governorate, user.Address, user.Birthdate); err != nil {
		println(err.Error())
		return err
	}
	return nil
}
