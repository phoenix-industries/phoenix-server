package apiservice

import (
	"net/http"
	"strconv"
	"time"

	"github.com/phoenix-industries/phoenix-server/pkg/auth"
	"github.com/phoenix-industries/phoenix-server/pkg/database/models"
	"github.com/phoenix-industries/phoenix-server/pkg/httputil"
)

type userUpdateData struct {
	Name        *string    `json:"name"`
	Email       *string    `json:"email"`
	Phone       *string    `json:"phone"`
	PictureID   *string    `json:"picture_id"`
	City        *string    `json:"city"`
	Governorate *string    `json:"governorate"`
	Address     *string    `json:"address"`
	Birthdate   *time.Time `json:"birthdate"`
}

func (s *Service) HandleGetUsers(w http.ResponseWriter, r *http.Request) *httputil.Response {
	role, err := httputil.GetUserRole(s.auth, r)
	if err != nil {
		return httputil.NewStatusError(err, "failed to get user role", http.StatusInternalServerError).Response()
	}
	if !role.Allowed(auth.RoleAdmin) {
		return httputil.NewStatusError(nil, "user not allowed", http.StatusForbidden).Response()
	}

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

	ctx := r.Context()
	users, err := models.UserGetAll(ctx, s.db.Pool(), limit, offset)
	if err != nil {
		return httputil.NewStatusError(err, "failed to get users", http.StatusInternalServerError).Response()
	}

	return httputil.NewResponseOK(http.StatusOK, users)
}

func (s *Service) HandleGetUserByID(w http.ResponseWriter, r *http.Request) *httputil.Response {
	id := r.PathValue("id")
	if id == "" {
		return httputil.NewStatusError(nil, "invalid request", http.StatusBadRequest).Response()
	}
	if id == "me" {
		cID, err := httputil.GetUserID(s.auth, r)
		if err != nil {
			return httputil.NewStatusError(err, "failed to get user id", http.StatusInternalServerError).Response()
		}
		id = cID
	}

	user, err := models.UserGetByID(r.Context(), s.db.Pool(), id)
	if err != nil {
		return httputil.NewStatusError(err, "failed to get user", http.StatusInternalServerError).Response()
	}
	if user == nil {
		return httputil.NewStatusError(nil, "user not found", http.StatusNotFound).Response()
	}

	return httputil.NewResponseOK(http.StatusOK, user)
}

func (s *Service) HandleUpdateUser(w http.ResponseWriter, r *http.Request) *httputil.Response {
	id := r.PathValue("id")
	if id == "" {
		return httputil.NewStatusError(nil, "invalid request", http.StatusBadRequest).Response()
	}
	if id == "me" {
		cID, err := httputil.GetUserID(s.auth, r)
		if err != nil {
			return httputil.NewStatusError(err, "failed to get user id", http.StatusInternalServerError).Response()
		}
		id = cID
	}

	var updateData userUpdateData
	if err := httputil.BodyJSON(w, r, &updateData); err != nil {
		return httputil.ErrInvalidBody.Response()
	}
	defer r.Body.Close()

	role, err := httputil.GetUserRole(s.auth, r)
	if err != nil {
		return httputil.NewStatusError(err, "failed to get user role", http.StatusInternalServerError).Response()
	}
	if !role.Allowed(auth.RoleAdmin) {
		cID, err := httputil.GetUserID(s.auth, r)
		if err != nil {
			return httputil.NewStatusError(err, "failed to get user id", http.StatusInternalServerError).Response()
		}
		if cID != id {
			return httputil.NewStatusError(nil, "user not allowed", http.StatusForbidden).Response()
		}
	}

	user, err := models.UserGetByID(r.Context(), s.db.Pool(), id)
	if err != nil {
		return httputil.NewStatusError(err, "failed to get user", http.StatusInternalServerError).Response()
	}
	if user == nil {
		return httputil.NewStatusError(nil, "user not found", http.StatusNotFound).Response()
	}
	if updateData.Name != nil {
		user.Name = *updateData.Name
	}
	if updateData.Email != nil {
		user.Email = *updateData.Email
	}
	if updateData.Phone != nil {
		user.Phone = *updateData.Phone
	}
	if updateData.PictureID != nil {
		user.PictureID = updateData.PictureID
	}
	if updateData.City != nil {
		user.City = updateData.City
	}
	if updateData.Governorate != nil {
		user.Governorate = updateData.Governorate
	}
	if updateData.Address != nil {
		user.Address = updateData.Address
	}
	if updateData.Birthdate != nil {
		user.Birthdate = *updateData.Birthdate
	}
	if err := models.UserUpdate(r.Context(), s.db.Pool(), user); err != nil {
		return httputil.NewStatusError(err, "failed to update user", http.StatusInternalServerError).Response()
	}

	return httputil.NewResponseOK(http.StatusNoContent, nil)
}
