package apiservice

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/phoenix-industries/phoenix-server/pkg/database/models"
	"github.com/phoenix-industries/phoenix-server/pkg/httputil"
)

type productCategoryData struct {
	Name *string `json:"name"`
}

func (s *Service) HandleCreateProductCategory(w http.ResponseWriter, r *http.Request) *httputil.Response {
	var data productCategoryData
	if err := httputil.BodyJSON(w, r, &data); err != nil {
		return httputil.ErrInvalidBody.Response()
	}
	if data.Name == nil || *data.Name == "" {
		return httputil.ErrBadRequest.Response()
	}

	id, err := s.auth.GenerateID()
	if err != nil {
		return httputil.NewStatusError(err, "failed to create product category", http.StatusInternalServerError).Response()
	}

	category := models.ProductCategory{
		Model: models.Model{ID: id},
		Name:  *data.Name,
	}
	if err = models.ProductCategoryInsert(r.Context(), s.db.Pool(), &category); err != nil {
		return httputil.NewStatusError(err, "failed to get category categories", http.StatusInternalServerError).Response()
	}

	return httputil.NewResponseOK(http.StatusCreated, category)
}

func (s *Service) HandleListProductCategories(w http.ResponseWriter, r *http.Request) *httputil.Response {
	categories, err := models.ProductCategoryGetAll(r.Context(), s.db.Pool())
	if err != nil {
		return httputil.NewStatusError(err, "failed to get category categories", http.StatusInternalServerError).Response()
	}
	return httputil.NewResponseOK(http.StatusOK, categories)
}

func (s *Service) HandleGetProductCategoryByID(w http.ResponseWriter, r *http.Request) *httputil.Response {
	id := r.PathValue("id")
	if id == "" {
		return httputil.ErrBadRequest.Response()
	}

	category, err := models.ProductCategoryGetByID(r.Context(), s.db.Pool(), id)
	if err != nil {
		return httputil.NewStatusError(err, "failed to get category category", http.StatusInternalServerError).Response()
	}

	return httputil.NewResponseOK(http.StatusOK, category)
}

func (s *Service) HandleUpdateProductCategory(w http.ResponseWriter, r *http.Request) *httputil.Response {
	id := r.PathValue("id")
	if id == "" {
		return httputil.ErrBadRequest.Response()
	}

	var updateData productCategoryData
	if err := httputil.BodyJSON(w, r, &updateData); err != nil {
		return httputil.ErrInvalidBody.Response()
	}

	ctx := r.Context()
	err := s.db.InTx(ctx, func(tx pgx.Tx) error {
		category, err := models.ProductCategoryGetByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get category category: %w", err)
		}
		if category == nil {
			return httputil.ErrNotFound
		}

		if updateData.Name != nil {
			category.Name = *updateData.Name
		}

		if err := models.ProductCategoryUpdate(ctx, tx, category); err != nil {
			return httputil.NewStatusError(err, "failed to update category", http.StatusInternalServerError)
		}

		return nil
	})
	if err != nil {
		return httputil.ResponseFromError(err)
	}

	return httputil.NewResponseOK(http.StatusNoContent, nil)
}

func (s *Service) HandleDeleteProductCategory(w http.ResponseWriter, r *http.Request) *httputil.Response {
	id := r.PathValue("id")
	if id == "" {
		return httputil.ErrBadRequest.Response()
	}

	ctx := r.Context()
	err := s.db.InTx(ctx, func(tx pgx.Tx) error {
		category, err := models.ProductCategoryGetByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get category: %w", err)
		}
		if category == nil {
			return httputil.NewStatusError(nil, "invalid category", http.StatusBadRequest)
		}

		if err := models.ProductCategoryDelete(r.Context(), s.db.Pool(), id); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return httputil.ResponseFromError(err)
	}

	return httputil.NewResponseOK(http.StatusNoContent, nil)
}
