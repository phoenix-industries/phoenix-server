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
	ImageID      *string    `db:"image_id" json:"image_id"`
	Name         string     `db:"name" json:"name"`
	Price        int64      `db:"price" json:"price"`
	Discount     int64      `db:"discount" json:"discount"`
	Quantity     int        `db:"quantity" json:"quantity"`
	CategoryID   string     `db:"category_id" json:"category_id"`
	Condition    string     `db:"condition" json:"condition"`
	MinimumAge   int        `db:"minimum_age" json:"minimum_age"`
	MaximumAge   int        `db:"maximum_age" json:"maximum_age"`
	TargetGender string     `db:"target_gender" json:"target_gender"`
	Description  string     `db:"description" json:"description"`
	Donated      bool       `db:"donated" json:"donated"`
	Approved     bool       `db:"approved" json:"approved"`
	Category     string     `db:"category" json:"category"`
	Brand        *string    `db:"brand" json:"brand"`
	Tags         *string    `db:"tags" json:"tags,omitempty"`
	ReviewedAt   *time.Time `db:"reviewed_at" json:"reviewed_at"`
}

func (p *Product) Validate() error {
	if p.UserID == "" {
		return errors.New("user id is required")
	}
	if p.CategoryID == "" {
		return errors.New("product category is required")
	}
	if p.Name == "" {
		return errors.New("product name is required")
	}
	if p.Condition == "" {
		return errors.New("product condition is required")
	}
	if p.Price < 0 {
		return errors.New("product price cannot be less than 0")
	}
	if p.Discount < 0 {
		return errors.New("product discount cannot be less than 0")
	}
	if p.Discount > p.Price {
		return errors.New("product discount cannot be greater than product price")
	}
	if p.MinimumAge < 3 {
		return errors.New("minimum age must be greater than or equal to 3")
	}
	if p.MaximumAge < 0 {
		return errors.New("maximum age must be greater than or equal to 0")
	}
	if p.MaximumAge > 0 && p.MaximumAge < p.MinimumAge {
		return errors.New("maximum age must be greater than or equal to minimum age")
	}
	if p.TargetGender == "" {
		return errors.New("target gender is required")
	}
	if p.Description == "" {
		return errors.New("product description is required")
	}
	if p.Donated && (p.Price != 0 || p.Discount != 0) {
		return errors.New("donated product cannot have price or discount")
	}
	return nil
}

func ProductInsert(ctx context.Context, db database.DB, product *Product) error {
	if err := product.Validate(); err != nil {
		return err
	}
	query := `
		INSERT INTO products
		(id, user_id, name, price, discount, quantity, category_id, condition, donated, minimum_age, maximum_age, target_gender, description, tags, image_id, approved)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, false)
	`
	_, err := db.Exec(
		ctx,
		query,
		product.ID,
		product.UserID,
		product.Name,
		product.Price,
		product.Discount,
		product.Quantity,
		product.CategoryID,
		product.Condition,
		product.Donated,
		product.MinimumAge,
		product.MaximumAge,
		product.TargetGender,
		product.Description,
		product.Tags,
		product.ImageID,
	)
	return err
}

func ProductGetByID(ctx context.Context, db database.DB, id string, requireApproval bool) (*Product, error) {
	query := `
		SELECT p.*, pc.name as category
		FROM products p
		LEFT JOIN product_categories pc
			ON p.category_id = pc.id AND pc.deleted_at IS null
		WHERE p.id = $1 AND p.deleted_at IS null
	`
	if requireApproval {
		query += " AND p.approved = true"
	}
	var product Product
	err := pgxscan.Get(ctx, db, &product, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

type ProductFilter struct {
	Name       string
	CategoryID string
	Condition  string
	Price      [2]int
}

type ProductListData struct {
	Product
	User User `db:"user" json:"user"`
}

func ProductList(ctx context.Context, db database.DB, requireApproval bool, limit int, offset int, filter *ProductFilter) ([]*Product, error) {
	args := []any{limit, offset}
	query := `
		SELECT
			p.*,
			pc.name as category
		FROM products p
		LEFT JOIN product_categories pc
			ON p.category_id = pc.id AND pc.deleted_at IS null
		WHERE p.deleted_at IS null
	`
	if requireApproval {
		query += " AND p.approved = true"
	}
	if filter != nil {
		if filter.Name != "" {
			args = append(args, filter.Name)
			query += fmt.Sprintf(" AND p.name ILIKE '%%$%d%%'", len(args))
		}
		if filter.CategoryID != "" {
			args = append(args, filter.CategoryID)
			query += fmt.Sprintf(" AND pc.id = $%d", len(args))
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
		GROUP BY p.id, pc.id
		ORDER BY p.created_at DESC
		LIMIT $1
		OFFSET $2
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

func ProductListByIDs(ctx context.Context, db database.DB, ids []string) ([]*Product, error) {
	query := `
		SELECT *
		FROM products
		WHERE id = ANY($1::text[])
	`
	var products []*Product
	if err := pgxscan.Select(ctx, db, &products, query, ids); err != nil {
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
		SET name = $2, price = $3, discount = $4, quantity = $5, category_id = $6, condition = $7, tags = $8,
			minimum_age = $9, maximum_age = $10, description = $11, image_id = $12,
			updated_at = CURRENT_TIMESTAMP, approved = false
		WHERE id = $1 AND deleted_at IS null
	`
	_, err := db.Exec(
		ctx,
		query,
		product.ID,
		product.Name,
		product.Price,
		product.Discount,
		product.Quantity,
		product.CategoryID,
		product.Condition,
		product.Tags,
		product.MinimumAge,
		product.MaximumAge,
		product.Description,
		product.ImageID,
	)
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

func ProductUpdateQuantityByID(ctx context.Context, db database.DB, id string, quantity int) error {
	query := `
		UPDATE products
		SET quantity = $2
		WHERE id = $1 AND deleted_at IS null
	`
	_, err := db.Exec(ctx, query, id, quantity)
	return err
}
