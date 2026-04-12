package apiservice

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/phoenix-industries/phoenix-server/pkg/auth"
	"github.com/phoenix-industries/phoenix-server/pkg/database/models"
	"github.com/phoenix-industries/phoenix-server/pkg/httputil"
)

func (s *Service) HandleListProducts(w http.ResponseWriter, r *http.Request) *httputil.Response {
	limit := 10
	offset := 0

	query := r.URL.Query()
	if limitQuery := query.Get("limit"); limitQuery != "" {
		limit, _ = strconv.Atoi(limitQuery)
		if limit == 0 {
			limit = 10
		}
	}
	if offsetQuery := query.Get("offset"); offsetQuery != "" {
		offset, _ = strconv.Atoi(offsetQuery)
	}

	if limit > 20 {
		return httputil.NewStatusError(nil, "limit cannot be greater than 20", http.StatusBadRequest).Response()
	}

	price := r.URL.Query().Get("price")
	priceSlice := strings.SplitN(price, "-", 2)
	var priceFilter [2]int
	if len(priceSlice) == 2 {
		priceFilter[0], _ = strconv.Atoi(priceSlice[0])
		if priceFilter[0] < 0 || priceFilter[0] >= math.MaxInt32-1 {
			priceFilter[0] = 0
		}
		priceFilter[1], _ = strconv.Atoi(priceSlice[1])
		if priceFilter[1] < 0 || priceFilter[1] >= math.MaxInt32-1 {
			priceFilter[1] = 0
		}
	}

	filter := models.ProductFilter{
		Name:      r.URL.Query().Get("filter"),
		Category:  r.URL.Query().Get("category"),
		Condition: r.URL.Query().Get("condition"),
		Price:     priceFilter,
	}

	role, err := httputil.GetUserRole(s.auth, r)
	if err != nil {
		return httputil.NewStatusError(err, "failed to get user role", http.StatusInternalServerError).Response()
	}
	requireApproval := !role.Allowed(auth.RoleManager)

	products, err := models.ProductList(r.Context(), s.db.Pool(), requireApproval, limit, offset, &filter)
	if err != nil {
		return httputil.NewStatusError(err, "failed to get products", http.StatusInternalServerError).Response()
	}

	return httputil.NewResponseOK(http.StatusOK, products)
}

func (s *Service) HandleGetProductByID(w http.ResponseWriter, r *http.Request) *httputil.Response {
	id := r.PathValue("id")
	if id == "" {
		return httputil.ErrBadRequest.Response()
	}

	role, err := httputil.GetUserRole(s.auth, r)
	if err != nil {
		return httputil.ResponseFromError(err)
	}
	requireApproval := !role.Allowed(auth.RoleManager)

	product, err := models.ProductGetByID(r.Context(), s.db.Pool(), id, requireApproval)
	if err != nil {
		return httputil.NewStatusError(err, "failed to get product", http.StatusInternalServerError).Response()
	}

	return httputil.NewResponseOK(http.StatusOK, product)
}

type productUpdateData struct {
	Name         *string `json:"name"`
	Price        *int    `json:"price"`
	CategoryID   *string `json:"category_id"`
	Condition    *string `json:"condition"`
	MinimumAge   *int    `json:"minimum_age"`
	MaximumAge   *int    `json:"maximum_age"`
	TargetGender *string `json:"target_gender"`
	Description  *string `json:"description"`
	Tags         *string `json:"tags"`
	Category     *string `json:"category"`
}

func (s *Service) HandleUpdateProduct(w http.ResponseWriter, r *http.Request) *httputil.Response {
	id := r.PathValue("id")
	if id == "" {
		return httputil.ErrBadRequest.Response()
	}

	var updateData productUpdateData
	if err := httputil.BodyJSON(r, &updateData); err != nil {
		return httputil.ErrInvalidBody.Response()
	}
	defer r.Body.Close()

	userID, err := httputil.GetUserID(s.auth, r)
	if err != nil {
		return httputil.ResponseFromError(err)
	}

	ctx := r.Context()
	err = s.db.InTx(ctx, func(tx pgx.Tx) error {
		user, err := models.UserGetByID(ctx, tx, userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		if user == nil {
			return httputil.ErrUnauthorized
		}

		product, err := models.ProductGetByID(ctx, tx, id, false)
		if err != nil {
			return fmt.Errorf("failed to get product: %w", err)
		}
		if product == nil {
			return httputil.ErrNotFound
		}
		if product.UserID != user.ID {
			return httputil.ErrUnauthorized
		}

		if updateData.Name != nil {
			product.Name = *updateData.Name
		}
		if updateData.Price != nil {
			product.Price = *updateData.Price
		}
		if updateData.CategoryID != nil {
			product.CategoryID = *updateData.CategoryID
		}
		if updateData.Condition != nil {
			product.Condition = *updateData.Condition
		}
		if updateData.MinimumAge != nil {
			product.MinimumAge = *updateData.MinimumAge
		}
		if updateData.MaximumAge != nil {
			product.MaximumAge = *updateData.MaximumAge
		}
		if updateData.TargetGender != nil {
			product.TargetGender = *updateData.TargetGender
		}
		if updateData.Description != nil {
			product.Description = *updateData.Description
		}
		if updateData.Tags != nil {
			product.Tags = updateData.Tags
		}

		if err := models.ProductUpdate(ctx, tx, product); err != nil {
			return httputil.NewStatusError(err, "failed to update product", http.StatusInternalServerError)
		}

		return nil
	})
	if err != nil {
		return httputil.ResponseFromError(err)
	}

	return httputil.NewResponseOK(http.StatusOK, nil)
}

func (s *Service) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) *httputil.Response {
	id := r.PathValue("id")
	if id == "" {
		return httputil.ErrBadRequest.Response()
	}

	userID, err := httputil.GetUserID(s.auth, r)
	if err != nil {
		return httputil.ResponseFromError(err)
	}

	ctx := r.Context()
	err = s.db.InTx(ctx, func(tx pgx.Tx) error {
		user, err := models.UserGetByID(ctx, tx, userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		if user == nil {
			return httputil.NewStatusError(nil, "invalid credentials", http.StatusUnauthorized)
		}

		product, err := models.ProductGetByID(ctx, tx, id, false)
		if err != nil {
			return fmt.Errorf("failed to get product: %w", err)
		}
		if product == nil {
			return httputil.NewStatusError(nil, "invalid product", http.StatusBadRequest)
		}
		if product.UserID != user.ID {
			return httputil.NewStatusError(nil, "invalid credentials", http.StatusUnauthorized)
		}

		if err := models.ProductDelete(r.Context(), s.db.Pool(), id); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return httputil.ResponseFromError(err)
	}

	return httputil.NewResponseOK(http.StatusOK, nil)
}
