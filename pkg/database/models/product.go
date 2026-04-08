package models

import (
	"context"
	"errors"
	"strconv"
	"strings"
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
	MiniumAge    int        `db:"minimum_age" json:"minimum_age"`
	MaximumAge   int        `db:"maximum_age" json:"maximum_age"`
	TargetGender string     `db:"target_gender" json:"target_gender"`
	Description  string     `db:"description" json:"description"`
	Tags         *string    `db:"tags" json:"tags"`
	ReviewedAt   *time.Time `db:"reviewed_at" json:"reviewed_at"`
	Approved     bool       `db:"approved" json:"approved"`
	Category     string     `db:"category" json:"category"`
}

func ProductGetByID(ctx context.Context, db database.DB, id string, requireApproval bool) (*Product, error) {
	stmt := `SELECT p.*, pc.name as category FROM products as p LEFT JOIN product_categories as pc ON p.category_id = pc.id AND pc.deleted_at IS NULL WHERE p.id = $1 AND p.deleted_at IS NULL AND p.approved = $2`
	var product Product
	err := pgxscan.Get(ctx, db, &product, stmt, id, requireApproval)
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
	var sb strings.Builder
	args := []any{requireApproval, limit, offset}
	sb.WriteString("SELECT p.*, pc.name as category FROM products as p LEFT JOIN product_categories as pc ON p.category_id = pc.id AND pc.deleted_at IS NULL WHERE p.deleted_at IS NULL AND p.approved = $1")
	if filter != nil {
		if filter.Name != "" {
			args = append(args, filter.Name)
			sb.WriteString(" AND p.name ILIKE '%$")
			sb.WriteString(strconv.Itoa(len(args)))
			sb.WriteString("%'")
		}
		if filter.Category != "" {
			args = append(args, filter.Category)
			sb.WriteString(" AND pc.name = $")
			sb.WriteString(strconv.Itoa(len(args)))
		}
		if filter.Condition != "" {
			args = append(args, filter.Condition)
			sb.WriteString(" AND p.condition = $")
			sb.WriteString(strconv.Itoa(len(args)))
		}
		if filter.Price[0] != 0 || filter.Price[1] != 0 {
			args = append(args, filter.Price[0], filter.Price[1])
			sb.WriteString(" AND p.price BETWEEN $")
			sb.WriteString(strconv.Itoa(len(args) - 1))
			sb.WriteString(" AND $")
			sb.WriteString(strconv.Itoa(len(args)))
		}
	}
	sb.WriteString(" ORDER BY p.created_at DESC LIMIT $2 OFFSET $3")
	stmt := sb.String()
	var products []*Product
	err := pgxscan.Select(ctx, db, &products, stmt, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return products, nil
}

func ProductUpdate(ctx context.Context, db database.DB, product *Product) error {
	stmt := `UPDATE products SET name = $2, price = $3, category_id = $4, condition = $5, tags = $6, minium_age = $7, maximum_age = $8, description = $9, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`
	_, err := db.Exec(ctx, stmt, product.ID, product.Name, product.Price, product.CategoryID, product.Condition, product.Tags, product.MiniumAge, product.MaximumAge, product.Description)
	return err
}

func ProductDelete(ctx context.Context, db database.DB, id string) error {
	stmt := `UPDATE products SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`
	_, err := db.Exec(ctx, stmt, id)
	return err
}
