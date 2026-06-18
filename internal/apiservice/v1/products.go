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

func (s *Service) HandleCreateProduct(w http.ResponseWriter, r *http.Request) *httputil.Response {
	var data models.Product
	if err := httputil.BodyJSON(w, r, &data); err != nil {
		return httputil.ErrInvalidBody.Response()
	}

	userID, err := httputil.GetUserID(s.auth, r)
	if err != nil {
		return httputil.ErrUnauthorized.Response()
	}

	ctx := r.Context()
	err = s.db.InTx(ctx, func(tx pgx.Tx) error {
		user, err := models.UserGetByID(ctx, tx, userID)
		if err != nil {
			return httputil.NewStatusError(err, "failed to get user", http.StatusInternalServerError)
		}
		if user == nil {
			return httputil.ErrUnauthorized
		}

		category, err := models.ProductCategoryGetByID(ctx, tx, data.CategoryID)
		if err != nil {
			return httputil.NewStatusError(err, "failed to get category", http.StatusInternalServerError)
		}
		if category == nil {
			return httputil.NewStatusError(nil, "category not found", http.StatusBadRequest)
		}

		id, err := s.auth.GenerateID()
		if err != nil {
			return httputil.NewStatusError(err, "failed to generate id", http.StatusInternalServerError)
		}
		data.ID = id
		data.UserID = user.ID
		data.Approved = false

		if err := data.Validate(); err != nil {
			return httputil.NewStatusError(nil, err.Error(), http.StatusBadRequest)
		}

		if err := models.ProductInsert(ctx, tx, &data); err != nil {
			return httputil.NewStatusError(err, "failed to create product", http.StatusInternalServerError)
		}

		return nil
	})
	if err != nil {
		return httputil.ResponseFromError(err)
	}

	return httputil.NewResponseOK(http.StatusCreated, nil)
}

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

	filter := models.ProductFilter{
		Name:       r.URL.Query().Get("query"),
		UserID:     r.URL.Query().Get("user"),
		CategoryID: r.URL.Query().Get("category"),
		Condition:  r.URL.Query().Get("condition"),
	}

	if filter.UserID == "me" {
		userID, err := httputil.GetUserID(s.auth, r)
		if err != nil {
			return httputil.NewStatusError(err, "failed to get user id", http.StatusUnauthorized).Response()
		}
		filter.UserID = userID
	}

	price := r.URL.Query().Get("price")
	if price != "" {
		filter.Price = &models.ProductFilterPrice{}
		priceSlice := strings.SplitN(price, "-", 2)
		if len(priceSlice) != 2 {
			return httputil.NewStatusError(nil, "invalid price range format", http.StatusBadRequest).Response()
		}
		priceMin, err := strconv.Atoi(priceSlice[0])
		if err != nil || priceMin < 0 || priceMin > math.MaxInt32 {
			return httputil.NewStatusError(err, "invalid price min range", http.StatusBadRequest).Response()
		}
		filter.Price.Min = priceMin
		priceMax, err := strconv.Atoi(priceSlice[1])
		if err != nil || priceMax < 0 || priceMax > math.MaxInt32 {
			return httputil.NewStatusError(err, "invalid price max range", http.StatusBadRequest).Response()
		}
		if priceMax != 0 && priceMax < priceMin {
			return httputil.NewStatusError(nil, "invalid price max range", http.StatusBadRequest).Response()
		}
		filter.Price.Max = priceMax
	}

	requireApproval := true
	if role, err := httputil.GetUserRole(s.auth, r); err == nil {
		requireApproval = !role.Allowed(auth.RoleManager)
	}

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

	requireApproval := true
	if role, err := httputil.GetUserRole(s.auth, r); err == nil {
		requireApproval = !role.Allowed(auth.RoleManager)
	}

	product, err := models.ProductGetByID(r.Context(), s.db.Pool(), id, requireApproval)
	if err != nil {
		return httputil.NewStatusError(err, "failed to get product", http.StatusInternalServerError).Response()
	}

	return httputil.NewResponseOK(http.StatusOK, product)
}

type productUpdateData struct {
	Name         *string `json:"name"`
	ImageID      *string `json:"image_id"`
	Price        *int64  `json:"price"`
	Discount     *int64  `json:"discount"`
	Quantity     *int    `json:"quantity"`
	CategoryID   *string `json:"category_id"`
	Condition    *string `json:"condition"`
	MinimumAge   *int    `json:"minimum_age"`
	MaximumAge   *int    `json:"maximum_age"`
	TargetGender *string `json:"target_gender"`
	Description  *string `json:"description"`
	Tags         *string `json:"tags"`
	Category     *string `json:"category"`
	Brand        *string `json:"brand"`
}

func (s *Service) HandleUpdateProduct(w http.ResponseWriter, r *http.Request) *httputil.Response {
	id := r.PathValue("id")
	if id == "" {
		return httputil.ErrBadRequest.Response()
	}

	var updateData productUpdateData
	if err := httputil.BodyJSON(w, r, &updateData); err != nil {
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
		if updateData.Discount != nil {
			product.Discount = *updateData.Discount
		}
		if updateData.Quantity != nil {
			product.Quantity = *updateData.Quantity
		}
		if updateData.CategoryID != nil {
			product.CategoryID = *updateData.CategoryID
		}
		if updateData.Condition != nil {
			product.Condition = *updateData.Condition
		}
		if updateData.MinimumAge != nil {
			product.MinimumAge = updateData.MinimumAge
		}
		if updateData.MaximumAge != nil {
			product.MaximumAge = updateData.MaximumAge
		}
		if updateData.TargetGender != nil {
			product.TargetGender = updateData.TargetGender
		}
		if updateData.Description != nil {
			product.Description = *updateData.Description
		}
		if updateData.Tags != nil {
			product.Tags = updateData.Tags
		}
		if updateData.ImageID != nil {
			product.ImageID = updateData.ImageID
		}
		if updateData.Brand != nil {
			product.Brand = updateData.Brand
		}

		if err := models.ProductUpdate(ctx, tx, product); err != nil {
			return httputil.NewStatusError(err, "failed to update product", http.StatusInternalServerError)
		}

		return nil
	})
	if err != nil {
		return httputil.ResponseFromError(err)
	}

	return httputil.NewResponseOK(http.StatusNoContent, nil)
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
			return httputil.ErrUnauthorized
		}

		product, err := models.ProductGetByID(ctx, tx, id, false)
		if err != nil {
			return fmt.Errorf("failed to get product: %w", err)
		}
		if product == nil {
			return httputil.ErrBadRequest
		}
		if product.UserID != user.ID {
			return httputil.ErrUnauthorized
		}

		if err := models.ProductDelete(r.Context(), s.db.Pool(), id); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return httputil.ResponseFromError(err)
	}

	return httputil.NewResponseOK(http.StatusNoContent, nil)
}

type productBuyData struct {
	Products []struct {
		ID       string `json:"id"`
		Quantity int    `json:"quantity"`
	} `json:"products"`
	ShippingInfo *struct {
		FullName string  `json:"full_name"`
		Phone    string  `json:"phone"`
		City     string  `json:"city"`
		Address  string  `json:"address"`
		Note     *string `json:"note"`
	} `json:"shipping_info"`
}

func (s *Service) HandleBuyProduct(w http.ResponseWriter, r *http.Request) *httputil.Response {
	const profitMin = 20 * 100
	const profitMax = 1000 * 100
	const profitPercent = 10 // %
	const shippingFee = 50 * 100
	const maxDonatedCount = 2
	const maxProductCount = 10

	var data productBuyData
	if err := httputil.BodyJSON(w, r, &data); err != nil {
		return httputil.ErrInvalidBody.Response()
	}
	if len(data.Products) == 0 || len(data.Products) > maxProductCount {
		return httputil.ErrBadRequest.Response()
	}
	fmt.Printf("data.Products: %+v\n", data.Products)
	fmt.Printf("data.shippingInfo: %+v\n", data.ShippingInfo)
	productIDs := make([]string, len(data.Products))
	quantityMap := make(map[string]int, len(data.Products))
	for i, p := range data.Products {
		if p.ID == "" || p.Quantity <= 0 {
			return httputil.ErrBadRequest.Response()
		}
		if _, ok := quantityMap[p.ID]; ok {
			return httputil.ErrBadRequest.Response()
		}
		productIDs[i] = p.ID
		quantityMap[p.ID] = p.Quantity
	}

	userID, err := httputil.GetUserID(s.auth, r)
	if err != nil {
		return httputil.ErrUnauthorized.Response()
	}

	ctx := r.Context()
	var invoice models.Invoice
	err = s.db.InTx(ctx, func(tx pgx.Tx) error {
		products, err := models.ProductListByIDs(ctx, tx, productIDs)
		if err != nil {
			return httputil.NewStatusError(err, "failed to get products", http.StatusInternalServerError)
		}
		if len(products) != len(productIDs) {
			return httputil.NewStatusError(nil, "invalid product ids", http.StatusBadRequest)
		}

		donatedCount, err := models.InvoiceItemGetDonatedCount(ctx, tx, userID)
		if err != nil {
			return httputil.NewStatusError(err, "failed to get donated count", http.StatusInternalServerError)
		}

		subtotal, discount := int64(0), int64(0)
		items := make([]models.InvoiceItem, len(products))
		for i, product := range products {
			quantity := quantityMap[product.ID]
			if quantity > product.Quantity {
				return httputil.NewStatusError(nil, "insufficient product quantity", http.StatusBadRequest)
			}
			quantity64 := int64(quantity)
			if product.Donated {
				if quantity > 1 || len(products) > 1 {
					return httputil.NewStatusError(nil, "donated product cannot have quantity greater than 1, don't be greedy pal", http.StatusBadRequest)
				}
				if donatedCount >= maxDonatedCount {
					return httputil.NewStatusError(nil, "donated count exceeds limit", http.StatusBadRequest)
				}
				donatedCount++
			}

			total := (product.Price - product.Discount) * quantity64
			if total < 0 {
				return httputil.NewStatusError(nil, "invalid product", http.StatusBadRequest)
			}
			subtotal += product.Price * quantity64
			discount += product.Discount * quantity64

			itemID, err := s.auth.GenerateID()
			if err != nil {
				return httputil.NewStatusError(err, "failed to generate id", http.StatusInternalServerError)
			}
			items[i] = models.InvoiceItem{
				Model: models.Model{
					ID: itemID,
				},
				ProductID: product.ID,
				Quantity:  quantity,
				Discount:  product.Discount,
				Price:     product.Price,
				Amount:    total,
			}

			if err := models.ProductUpdateQuantityByID(ctx, tx, product.ID, product.Quantity-quantity); err != nil {
				return httputil.NewStatusError(err, "failed to update product quantity", http.StatusInternalServerError)
			}

			if !product.Donated {
				sellerWallet, err := models.WalletGetByUserID(ctx, tx, product.UserID)
				if err != nil {
					return httputil.NewStatusError(err, "failed to get seller's wallet", http.StatusInternalServerError)
				}
				if sellerWallet == nil {
					return httputil.NewStatusError(nil, "invalid seller's wallet", http.StatusInternalServerError)
				}
				transactionID, err := s.auth.GenerateID()
				if err != nil {
					return httputil.NewStatusError(err, "failed to generate id", http.StatusInternalServerError)
				}
				profit := min(max(total/profitPercent, profitMin), profitMax)
				income := total - profit
				transaction := models.Transaction{
					Model: models.Model{
						ID: transactionID,
					},
					UserID:      sellerWallet.UserID,
					WalletID:    sellerWallet.ID,
					Debit:       0,
					Credit:      income,
					Currency:    "egp",
					Description: "sale",
					Status:      models.TransactionStatusSuccess,
				}
				if err := models.TransactionInsert(ctx, tx, &transaction); err != nil {
					return httputil.NewStatusError(err, "failed to create transaction", http.StatusInternalServerError)
				}
				if err := models.WalletTopupByID(ctx, tx, sellerWallet.ID, income); err != nil {
					return httputil.NewStatusError(err, "failed to update seller's wallet balance", http.StatusInternalServerError)
				}
			}
		}
		amount := subtotal - discount
		if data.ShippingInfo != nil {
			amount += shippingFee
		}

		wallet, err := models.WalletGetByUserID(ctx, tx, userID)
		if err != nil {
			return httputil.NewStatusError(err, "failed to get user's wallet", http.StatusInternalServerError)
		}
		if wallet == nil {
			return httputil.NewStatusError(nil, "invalid wallet", http.StatusInternalServerError)
		}
		if wallet.Balance < amount {
			return httputil.NewStatusError(nil, "insufficient funds", http.StatusBadRequest)
		}
		if err := models.WalletWithdrawByID(ctx, tx, wallet.ID, amount); err != nil {
			return httputil.NewStatusError(err, "failed to update user's wallet balance", http.StatusInternalServerError)
		}

		transactionID, err := s.auth.GenerateID()
		if err != nil {
			return httputil.NewStatusError(err, "failed to generate id", http.StatusInternalServerError)
		}
		transaction := models.Transaction{
			Model: models.Model{
				ID: transactionID,
			},
			UserID:      userID,
			WalletID:    wallet.ID,
			Debit:       amount,
			Credit:      0,
			Currency:    "egp",
			Description: "purchase",
			Status:      models.TransactionStatusSuccess,
		}
		if err := models.TransactionInsert(ctx, tx, &transaction); err != nil {
			return httputil.NewStatusError(err, "failed to create transaction", http.StatusInternalServerError)
		}

		invoiceID, err := s.auth.GenerateID()
		if err != nil {
			return httputil.NewStatusError(err, "failed to generate id", http.StatusInternalServerError)
		}
		invoice = models.Invoice{
			Model: models.Model{
				ID: invoiceID,
			},
			UserID:        userID,
			TransactionID: transaction.ID,
			Subtotal:      subtotal,
			Discount:      discount,
			Amount:        amount,
			Currency:      "egp",
			Description:   "purchase",
		}
		if err := models.InvoiceInsert(ctx, tx, &invoice); err != nil {
			return httputil.NewStatusError(err, "failed to create invoice", http.StatusInternalServerError)
		}

		for i := range items {
			items[i].InvoiceID = invoiceID
		}
		if err := models.InvoiceItemInsertBatch(ctx, tx, items); err != nil {
			return httputil.NewStatusError(err, "failed to create invoice item", http.StatusInternalServerError)
		}

		if data.ShippingInfo != nil {
			shippingID, err := s.auth.GenerateID()
			if err != nil {
				return httputil.NewStatusError(err, "failed to generate id", http.StatusInternalServerError)
			}
			shipping := models.Shipping{
				Model: models.Model{
					ID: shippingID,
				},
				UserID:    userID,
				InvoiceID: invoiceID,
				Fee:       shippingFee,
				FullName:  data.ShippingInfo.FullName,
				Phone:     data.ShippingInfo.Phone,
				City:      data.ShippingInfo.City,
				Address:   data.ShippingInfo.Address,
				Note:      data.ShippingInfo.Note,
			}
			if err := shipping.Validate(); err != nil {
				return httputil.NewStatusError(nil, "invalid shipping info: "+err.Error(), http.StatusBadRequest)
			}
			if err := models.ShippingInsert(ctx, tx, &shipping); err != nil {
				return httputil.NewStatusError(err, "failed to create shipping", http.StatusInternalServerError)
			}
		}

		return nil
	})
	if err != nil {
		return httputil.ResponseFromError(err)
	}

	return httputil.NewResponseOK(http.StatusOK, &invoice)
}
