package models

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/phoenix-industries/phoenix-server/pkg/database"
)

type Wallet struct {
	Model
	UserID   string `db:"user_id"`
	Name     string `db:"name"`
	Balance  int64  `db:"balance"`
	Currency string `db:"currency"`
}

func WalletInsert(ctx context.Context, db database.DB, wallet *Wallet) error {
	query := `
		INSERT INTO wallets
		(id, user_id, name, balance, currency)
		VALUES
		($1, $2, $3, $4, $5)
	`
	_, err := db.Exec(ctx, query, wallet.ID, wallet.UserID, wallet.Name, wallet.Balance, wallet.Currency)
	return err
}

func WalletGetByUserID(ctx context.Context, db database.DB, userID string) (*Wallet, error) {
	query := `
		SELECT *
		FROM wallets
		WHERE user_id = $1 AND deleted_at IS NULL
	`
	var wallet Wallet
	if err := pgxscan.Get(ctx, db, &wallet, query, userID); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &wallet, nil
}

func WalletUpdateBalanceByID(ctx context.Context, db database.DB, id string, balance int64) error {
	query := `
		UPDATE wallets
		SET balance = $2
		WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := db.Exec(ctx, query, id, balance)
	return err
}

func WalletUpdateBalanceByUserID(ctx context.Context, db database.DB, userID string, balance int64) error {
	query := `
		UPDATE wallets
		SET balance = $2
		WHERE user_id = $1 AND deleted_at IS NULL
	`
	_, err := db.Exec(ctx, query, userID, balance)
	return err
}

func WalletTopupByID(ctx context.Context, db database.DB, userID string, amount int64) error {
	query := `
		UPDATE wallets
		SET balance = balance + $2
		WHERE user_id = $1 AND deleted_at IS NULL
	`
	_, err := db.Exec(ctx, query, userID, amount)
	return err
}

func WalletWithdrawByID(ctx context.Context, db database.DB, userID string, amount int64) error {
	query := `
		UPDATE wallets
		SET balance = balance - $2
		WHERE user_id = $1 AND deleted_at IS NULL
	`
	_, err := db.Exec(ctx, query, userID, amount)
	return err
}
