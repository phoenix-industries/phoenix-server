package models

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/phoenix-industries/phoenix-server/pkg/database"
)

type Invoice struct {
	Model
	UserID        string `db:"user_id" json:"user_id"`
	TransactionID string `db:"transaction_id" json:"transaction_id"`
	Subtotal      int64  `db:"subtotal" json:"subtotal"`
	Discount      int64  `db:"discount" json:"discount"`
	Amount        int64  `db:"amount" json:"amount"`
	Currency      string `db:"currency" json:"currency"`
	Description   string `db:"description" json:"description"`
}

func InvoiceInsert(ctx context.Context, db database.DB, invoice *Invoice) error {
	query := `
		INSERT INTO invoices
		(id, user_id, transaction_id, subtotal, discount, amount, currency, description)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := db.Exec(ctx,
		query,
		invoice.ID,
		invoice.UserID,
		invoice.TransactionID,
		invoice.Subtotal,
		invoice.Discount,
		invoice.Amount,
		invoice.Currency,
		invoice.Description,
	)
	return err
}

func InvoiceGetByID(ctx context.Context, db database.DB, id string) (*Invoice, error) {
	query := `
		SELECT *
		FROM invoices
		WHERE id = $1 AND deleted_at IS null
	`
	var invoice Invoice
	if err := pgxscan.Get(ctx, db, &invoice, query, id); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &invoice, nil
}

type InvoiceItemWithProduct struct {
	InvoiceItem
	Product *Product `db:"product" json:"product"`
}

type InvoiceListData struct {
	Invoice
	Items    []*InvoiceItemWithProduct `db:"items" json:"items"`
	Shipping *Shipping                 `db:"shipping" json:"shipping,omitempty"`
}

func InvoiceListByUserID(ctx context.Context, db database.DB, userID string, limit int, offset int) ([]*InvoiceListData, error) {
	query := `
		SELECT
			i.id,
			i.created_at,
			i.updated_at,
			i.deleted_at,
			i.user_id,
			i.transaction_id,
			i.subtotal,
			i.discount,
			i.amount,
			i.currency,
			i.description,
			COALESCE(
				json_agg(
					json_build_object(
						'id',         ii.id,
						'quantity',   ii.quantity,
						'price',      ii.price,
						'product',    json_build_object(
							'id',            p.id,
							'name',          p.name,
							'price',         p.price,
							'discount',      p.discount,
							'quantity',      p.quantity,
							'category_id',   p.category_id,
							'condition',     p.condition,
							'minimum_age',   p.minimum_age,
							'maximum_age',   p.maximum_age,
							'target_gender', p.target_gender,
							'description',   p.description,
							'approved',      p.approved,
							'category',      p.category,
							'tags',          p.tags,
							'reviewed_at',   p.reviewed_at
						)
					)
				) FILTER (WHERE ii.id IS NOT NULL AND ii.deleted_at IS NULL),
				'[]'
			) AS items,
			COALESCE(
				json_build_object(
					'id',         s.id,
					'fee',        s.fee,
					'address',    s.address,
					'note',       s.note
				) FILTER (WHERE s.id IS NOT NULL AND s.deleted_at IS NULL),
				NULL
			) AS shipping
		FROM invoices i
		LEFT JOIN invoice_items ii ON ii.invoice_id = i.id
		LEFT JOIN products p ON p.id = ii.product_id
		LEFT JOIN shippings s ON s.invoice_id = i.id
		WHERE i.user_id = $1 AND i.deleted_at IS NULL
		GROUP BY i.id
		ORDER BY i.created_at DESC
		LIMIT $2
		OFFSET $3
	`
	var invoices []*InvoiceListData
	if err := pgxscan.Select(ctx, db, &invoices, query, userID, limit, offset); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return invoices, nil
}
