package models

import (
	"context"
	"errors"

	"github.com/phoenix-industries/phoenix-server/pkg/database"
)

type Shipping struct {
	Model
	UserID    string  `db:"user_id" json:"user_id"`
	InvoiceID string  `db:"invoice_id" json:"invoice_id"`
	Fee       int64   `db:"fee" json:"fee"`
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
		(id, user_id, invoice_id, fee, address, note)
		VALUES
		($1, $2, $3, $4, $5, $6)
	`
	_, err := db.Exec(
		ctx,
		query,
		shipping.ID,
		shipping.UserID,
		shipping.InvoiceID,
		shipping.Fee,
		shipping.Address,
		shipping.Note,
	)
	return err
}
