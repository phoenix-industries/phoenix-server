package apiservice

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/phoenix-industries/phoenix-server/pkg/auth"
	"github.com/phoenix-industries/phoenix-server/pkg/database/models"
	"github.com/phoenix-industries/phoenix-server/pkg/httputil"
)

func (s *Service) HandleListProducts(w http.ResponseWriter, r *http.Request) {
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
		httputil.Error(nil, "limit cannot be greater than 20", http.StatusBadRequest).WriteJSON(w)
		return
	}

	if limit > 20 {
		httputil.Error(nil, "limit cannot be greater than 20", http.StatusBadRequest).WriteJSON(w)
		return
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
		httputil.Error(err, "failed to get user role", http.StatusInternalServerError).WriteJSON(w)
		return
	}
	requireApproval := !role.Allowed(auth.RoleManager)

	products, err := models.ProductList(r.Context(), s.db.Pool(), requireApproval, limit, offset, &filter)
	if err != nil {
		println(err.Error())
		httputil.Error(err, "failed to get products", http.StatusInternalServerError).WriteJSON(w)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		httputil.Error(err, "failed to encode response", http.StatusInternalServerError).WriteJSON(w)
	}
}

func (s *Service) HandleGetProductByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		httputil.Error(nil, "invalid request", http.StatusBadRequest).WriteJSON(w)
		return
	}

	role, err := httputil.GetUserRole(s.auth, r)
	if err != nil {
		fmt.Println(err.Error())
		httputil.Error(err, "failed to get user role", http.StatusInternalServerError).WriteJSON(w)
		return
	}
	requireApproval := !role.Allowed(auth.RoleManager)

	product, err := models.ProductGetByID(r.Context(), s.db.Pool(), id, requireApproval)
	if err != nil {
		httputil.Error(err, "failed to get product", http.StatusInternalServerError).WriteJSON(w)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		httputil.Error(err, "failed to encode response", http.StatusInternalServerError).WriteJSON(w)
	}
}

func (s *Service) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		httputil.Error(nil, "invalid request", http.StatusBadRequest).WriteJSON(w)
		return
	}

	if err := models.ProductDelete(r.Context(), s.db.Pool(), id); err != nil {
		httputil.Error(err, "failed to delete product", http.StatusInternalServerError).WriteJSON(w)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "ok"}); err != nil {
		httputil.Error(err, "failed to encode response", http.StatusInternalServerError).WriteJSON(w)
	}
}
