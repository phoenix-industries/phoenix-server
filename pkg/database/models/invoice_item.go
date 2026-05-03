package models

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/phoenix-industries/phoenix-server/pkg/database"
)

type InvoiceItem struct {
	Model
	InvoiceID string `db:"invoice_id" json:"invoice_id"`
	ProductID string `db:"product_id" json:"product_id"`
	Quantity  int    `db:"quantity" json:"quantity"`
	Price     int64  `db:"price" json:"price"`
	Discount  int64  `db:"discount" json:"discount"`
	Amount    int64  `db:"amount" json:"amount"`
}

func InvoiceItemInsert(ctx context.Context, db database.DB, invoiceItem *InvoiceItem) error {
	query := `
		INSERT INTO invoice_items
		(id, invoice_id, product_id, quantity, price, discount, amount)
		VALUES
		($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := db.Exec(ctx,
		query,
		invoiceItem.ID,
		invoiceItem.InvoiceID,
		invoiceItem.ProductID,
		invoiceItem.Quantity,
		invoiceItem.Price,
		invoiceItem.Discount,
		invoiceItem.Amount,
	)
	return err
}

func InvoiceItemInsertBatch(ctx context.Context, db database.DB, items []InvoiceItem) error {
	_, err := db.CopyFrom(
		ctx,
		pgx.Identifier{"invoice_items"},
		[]string{"id", "invoice_id", "product_id", "quantity", "price", "discount", "amount"},
		pgx.CopyFromSlice(len(items), func(i int) ([]any, error) {
			return []any{
				items[i].ID,
				items[i].InvoiceID,
				items[i].ProductID,
				items[i].Quantity,
				items[i].Price,
				items[i].Discount,
				items[i].Amount,
			}, nil
		}),
	)
	return err
}

func InvoiceItemGetDonatedCount(ctx context.Context, db database.DB, userID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM products p
		INNER JOIN invoice_items ii ON ii.product_id = p.id
		INNER JOIN invoices i ON i.id = ii.invoice_id
		WHERE i.user_id = $1
			AND p.donated = true
			AND i.deleted_at IS NULL
			AND ii.deleted_at IS NULL
			AND i.created_at >= date_trunc('month', now())
			AND i.created_at < date_trunc('month', now()) + interval '1 month'
	`
	var count int
	if err := pgxscan.Get(ctx, db, &count, query, userID); err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}
