package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/phoenix-industries/phoenix-server/pkg/database"
)

type Product struct {
	Model
	UserID       string     `db:"user_id" json:"user_id"`
	Name         string     `db:"name" json:"name"`
	Price        int        `db:"price" json:"price"`
	CategoryID   string     `db:"category_id" json:"-"`
	Condition    string     `db:"condition" json:"condition"`
	MinimumAge   int        `db:"minimum_age" json:"minimum_age"`
	MaximumAge   int        `db:"maximum_age" json:"maximum_age"`
	TargetGender string     `db:"target_gender" json:"target_gender"`
	Description  string     `db:"description" json:"description"`
	Tags         *string    `db:"tags" json:"tags"`
	ReviewedAt   *time.Time `db:"reviewed_at" json:"reviewed_at"`
	Approved     bool       `db:"approved" json:"approved"`
	Category     string     `db:"category" json:"category"`
}

func ProductGetByID(ctx context.Context, db database.DB, id string, requireApproval bool) (*Product, error) {
	query := `
		SELECT p.*, pc.name as category
		FROM products p
		LEFT JOIN product_categories pc
			ON p.category_id = pc.id AND pc.deleted_at IS null
		WHERE p.id = $1 AND p.deleted_at IS null AND p.approved = $2
	`
	var product Product
	err := pgxscan.Get(ctx, db, &product, query, id, requireApproval)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

type ProductFilter struct {
	Name      string
	Category  string
	Condition string
	Price     [2]int
}

func ProductList(ctx context.Context, db database.DB, requireApproval bool, limit int, offset int, filter *ProductFilter) ([]*Product, error) {
	args := []any{requireApproval, limit, offset}
	query := `
		SELECT p.*, pc.name as category
		FROM products p
		LEFT JOIN product_categories pc
			ON p.category_id = pc.id AND pc.deleted_at IS null
		WHERE p.deleted_at IS null AND p.approved = $1
	`
	if filter != nil {
		if filter.Name != "" {
			args = append(args, filter.Name)
			query += fmt.Sprintf(" AND p.name ILIKE '%%$%d%%'", len(args))
		}
		if filter.Category != "" {
			args = append(args, filter.Category)
			query += fmt.Sprintf(" AND pc.name = $%d", len(args))
		}
		if filter.Condition != "" {
			args = append(args, filter.Condition)
			query += fmt.Sprintf(" AND p.condition = $%d", len(args))
		}
		if len(filter.Price) == 2 && (filter.Price[0] != 0 || filter.Price[1] != 0) {
			args = append(args, filter.Price[0], filter.Price[1])
			query += fmt.Sprintf(" AND p.price BETWEEN $%d AND $%d", len(args)-1, len(args))
		}
	}
	query += `
		ORDER BY p.created_at DESC
		LIMIT $2
		OFFSET $3
	`
	var products []*Product
	err := pgxscan.Select(ctx, db, &products, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return products, nil
}

func ProductUpdate(ctx context.Context, db database.DB, product *Product) error {
	query := `
		UPDATE products
		SET name = $2, price = $3, category_id = $4, condition = $5, tags = $6,
			minimum_age = $7, maximum_age = $8, description = $9,
			updated_at = CURRENT_TIMESTAMP, approved = false
		WHERE id = $1 AND deleted_at IS null
	`
	_, err := db.Exec(ctx, query, product.ID, product.Name, product.Price, product.CategoryID, product.Condition, product.Tags, product.MinimumAge, product.MaximumAge, product.Description)
	return err
}

func ProductDelete(ctx context.Context, db database.DB, id string) error {
	query := `
		UPDATE products
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS null`
	_, err := db.Exec(ctx, query, id)
	return err
}
