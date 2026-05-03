package models

import (
	"context"
	"errors"
	"time"

	"github.com/phoenix-industries/phoenix-server/pkg/database"
)

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusSuccess TransactionStatus = "success"
	TransactionStatusFailed  TransactionStatus = "failed"
)

type Transaction struct {
	Model
	UserID      string            `db:"user_id" json:"user_id"`
	WalletID    string            `db:"wallet_id" json:"wallet_id"`
	Debit       int64             `db:"debit" json:"debit"`
	Credit      int64             `db:"credit" json:"credit"`
	Currency    string            `db:"currency" json:"currency"`
	Description string            `db:"description" json:"description"`
	Status      TransactionStatus `db:"status" json:"status"`
	CreatedAt   time.Time         `db:"created_at" json:"created_at"`
}

func (t *Transaction) Validate() error {
	if t.UserID == "" {
		return errors.New("user_id is required")
	}
	if t.WalletID == "" {
		return errors.New("wallet_id is required")
	}
	if t.Currency == "" {
		return errors.New("currency is required")
	}
	if t.Description == "" {
		return errors.New("description is required")
	}
	if t.Status == "" {
		return errors.New("status is required")
	}
	return nil
}

func TransactionInsert(ctx context.Context, db database.DB, transaction *Transaction) error {
	if transaction == nil {
		return errors.New("transaction is nil")
	}
	if transaction.ID == "" {
		return errors.New("id is not set")
	}
	if err := transaction.Validate(); err != nil {
		return err
	}
	query := `
		INSERT INTO transactions
		(id, user_id, wallet_id, debit, credit, currency, description, status, created_at)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := db.Exec(ctx, query, transaction.ID, transaction.UserID, transaction.WalletID, transaction.Debit, transaction.Credit, transaction.Currency, transaction.Description, transaction.Status, transaction.CreatedAt)
	return err
}
