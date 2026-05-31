package models

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/phoenix-industries/phoenix-server/pkg/database"
)

type ProductCategory struct {
	Model
	Name string `db:"name" json:"name"`
}

func (c *ProductCategory) Validate() error {
	if c.Name == "" {
		return errors.New("category name is required")
	}
	return nil
}

func ProductCategoryInsert(ctx context.Context, db database.DB, category *ProductCategory) error {
	if err := category.Validate(); err != nil {
		return err
	}
	query := `
		INSERT INTO product_categories
		(id, name)
		VALUES
		($1, $2)
	`
	_, err := db.Exec(ctx, query, category.ID, category.Name)
	return err
}

func ProductCategoryGetByID(ctx context.Context, db database.DB, id string) (*ProductCategory, error) {
	query := `
		SELECT *
		FROM product_categories
		WHERE id = $1 AND deleted_at IS null
	`
	var category ProductCategory
	err := pgxscan.Get(ctx, db, &category, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

func ProductCategoryGetAll(ctx context.Context, db database.DB) ([]*ProductCategory, error) {
	query := `
		SELECT *
		FROM product_categories
		WHERE deleted_at IS null
		ORDER BY created_at
	`
	var categories []*ProductCategory
	err := pgxscan.Select(ctx, db, &categories, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return categories, nil
}

func ProductCategoryUpdate(ctx context.Context, db database.DB, category *ProductCategory) error {
	query := `
		UPDATE product_categories
		SET name = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS null
	`
	_, err := db.Exec(ctx, query, category.ID, category.Name)
	return err
}

func ProductCategoryDelete(ctx context.Context, db database.DB, id string) error {
	query := `
		UPDATE product_categories
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS null`
	_, err := db.Exec(ctx, query, id)
	return err
}
