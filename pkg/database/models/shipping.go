package models

import (
	"context"
	"errors"

	"github.com/phoenix-industries/phoenix-server/pkg/database"
	"github.com/phoenix-industries/phoenix-server/pkg/validate"
)

type Shipping struct {
	Model
	UserID    string  `db:"user_id" json:"user_id"`
	InvoiceID string  `db:"invoice_id" json:"invoice_id"`
	Fee       int64   `db:"fee" json:"fee"`
	FullName  string  `db:"full_name" json:"full_name"`
	Phone     string  `db:"phone" json:"phone"`
	City      string  `db:"city" json:"city"`
	Address   string  `db:"address" json:"address"`
	Note      *string `db:"note" json:"note,omitempty"`
}

func (s *Shipping) Validate() error {
	if s.ID == "" {
		return errors.New("id is required")
	}
	if s.UserID == "" {
		return errors.New("user id is required")
	}
	if s.InvoiceID == "" {
		return errors.New("invoice id is required")
	}
	if s.Fee <= 0 {
		return errors.New("fee must be greater than 0")
	}
	if s.FullName == "" {
		return errors.New("full name is required")
	}
	if s.Phone == "" {
		return errors.New("phone is required")
	}
	if err := validate.PhoneNumber(s.Phone); err != nil {
		return err
	}
	if s.City == "" {
		return errors.New("city is required")
	}
	if s.Address == "" {
		return errors.New("address is required")
	}
	return nil
}

func ShippingInsert(ctx context.Context, db database.DB, shipping *Shipping) error {
	if err := shipping.Validate(); err != nil {
		return err
	}
	query := `
		INSERT INTO shippings
		(id, user_id, invoice_id, fee, full_name, phone, city, address, note)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := db.Exec(
		ctx,
		query,
		shipping.ID,
		shipping.UserID,
		shipping.InvoiceID,
		shipping.Fee,
		shipping.FullName,
		shipping.Phone,
		shipping.City,
		shipping.Address,
		shipping.Note,
	)
	return err
}
