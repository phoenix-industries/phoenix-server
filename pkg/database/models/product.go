package models

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/phoenix-industries/phoenix-server/pkg/database"
)

type Product struct {
	Model
	UserID       string     `db:"user_id" json:"user_id,omitempty"`
	ImageID      *string    `db:"image_id" json:"image_id"`
	Name         string     `db:"name" json:"name"`
	Price        int64      `db:"price" json:"price"`
	Discount     int64      `db:"discount" json:"discount"`
	Quantity     int        `db:"quantity" json:"quantity"`
	CategoryID   string     `db:"category_id" json:"category_id,omitempty"`
	Condition    string     `db:"condition" json:"condition"`
	MinimumAge   *int       `db:"minimum_age" json:"minimum_age,omitempty"`
	MaximumAge   *int       `db:"maximum_age" json:"maximum_age,omitempty"`
	TargetGender *string    `db:"target_gender" json:"target_gender,omitempty"`
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
	if p.Price > math.MaxInt32 {
		return errors.New("product price too large")
	}
	if p.Discount < 0 {
		return errors.New("product discount cannot be less than 0")
	}
	if p.Discount > math.MaxInt32 {
		return errors.New("product discount too large")
	}
	if p.Discount > p.Price {
		return errors.New("product discount cannot be greater than product price")
	}
	if p.Donated && (p.Price > 0 || p.Discount > 0) {
		return errors.New("a donated product cannot have a price nor a discount")
	}
	if p.MinimumAge != nil && *p.MinimumAge < 3 {
		return errors.New("minimum age must be greater than or equal to 3")
	}
	if p.MaximumAge != nil && *p.MaximumAge < 0 {
		return errors.New("maximum age must be greater than or equal to 0")
	}
	if p.MaximumAge != nil && *p.MaximumAge > 0 && *p.MaximumAge < *p.MinimumAge {
		return errors.New("maximum age must be greater than or equal to minimum age")
	}
	if p.TargetGender != nil && *p.TargetGender == "" {
		return errors.New("target gender cannot be empty")
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
	query := `
		INSERT INTO products
		(id, user_id, name, price, discount, quantity, category_id, condition, donated, minimum_age, maximum_age, target_gender, description, tags, image_id, brand, approved)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, false)
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
		product.Brand,
	)
	return err
}

func ProductGetByID(ctx context.Context, db database.DB, id string, requireApproval bool) (*Product, error) {
	query := `
		SELECT p.*, pc.name as category
		FROM products p
		LEFT JOIN product_categories pc
			ON p.category_id = pc.id AND pc.deleted_at IS null
		WHERE p.id = $1 AND p.user_id IS NOT null AND p.deleted_at IS null
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

type ProductFilterPrice struct {
	Min int
	Max int
}

type ProductFilter struct {
	Name       string
	CategoryID string
	Condition  string
	Price      *ProductFilterPrice
}

type ProductListData struct {
	Product
	Category struct {
		ID   string `db:"id" json:"id"`
		Name string `db:"name" json:"name"`
	} `db:"category" json:"category"`
	User struct {
		ID          string  `db:"id" json:"id"`
		Name        string  `db:"name" json:"name"`
		PictureID   *string `db:"picture_id" json:"picture_id,omitempty"`
		City        *string `db:"city" json:"city,omitempty"`
		Governorate *string `db:"governorate" json:"governorate,omitempty"`
		Address     *string `db:"address" json:"address,omitempty"`
	} `db:"user" json:"user"`
}

func ProductList(ctx context.Context, db database.DB, requireApproval bool, limit int, offset int, filter *ProductFilter) ([]*ProductListData, error) {
	args := []any{limit, offset}
	query := `
		SELECT
			p.id,
			p.name,
			p.price,
			p.discount,
			p.quantity,
			p.condition,
			p.description,
			p.image_id,
			p.brand,
			p.tags,
			p.reviewed_at,
			p.created_at,
			p.updated_at,
			p.approved,
			json_build_object(
				'id', pc.id,
				'name', pc.name
			) as category,
			json_build_object(
				'id', u.id,
				'name', u.name,
				'picture_id', u.picture_id,
				'city', u.city,
				'governorate', u.governorate,
				'address', u.address
			) as "user"
		FROM products p
		LEFT JOIN users u
			ON p.user_id = u.id
		LEFT JOIN product_categories pc
			ON p.category_id = pc.id AND pc.deleted_at IS null
		WHERE p.user_id IS NOT null AND p.deleted_at IS null
	`
	if requireApproval {
		query += " AND p.approved = true"
	}
	if filter != nil {
		if filter.Name != "" {
			args = append(args, fmt.Sprintf("%%%s%%", filter.Name))
			query += fmt.Sprintf(" AND p.name ILIKE $%d", len(args))
		}
		if filter.CategoryID != "" {
			args = append(args, filter.CategoryID)
			query += fmt.Sprintf(" AND pc.id = $%d", len(args))
		}
		if filter.Condition != "" {
			args = append(args, filter.Condition)
			query += fmt.Sprintf(" AND p.condition = $%d", len(args))
		}
		if filter.Price != nil {
			if filter.Price.Min == 0 && filter.Price.Max == 0 {
				query += " AND p.price = 0"
			} else {
				if filter.Price.Min > 0 {
					args = append(args, filter.Price.Min)
					query += fmt.Sprintf(" AND p.price >= $%d", len(args))
				}
				if filter.Price.Max > 0 {
					args = append(args, filter.Price.Max)
					query += fmt.Sprintf(" AND p.price <= $%d", len(args))
				}
			}
		}
	}
	query += `
		GROUP BY p.id, pc.id, u.id
		ORDER BY p.created_at DESC
		LIMIT $1
		OFFSET $2
	`
	var products []*ProductListData
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
		WHERE id = ANY($1::text[]) AND user_id IS NOT null
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
			minimum_age = $9, maximum_age = $10, description = $11, image_id = $12, brand = $13,
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
		product.Brand,
	)
	return err
}

func ProductDelete(ctx context.Context, db database.DB, id string) error {
	query := `
		UPDATE products
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS null
	`
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
