package authservice

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/phoenix-industries/phoenix-server/pkg/auth"
	"github.com/phoenix-industries/phoenix-server/pkg/database/models"
	"github.com/phoenix-industries/phoenix-server/pkg/httputil"
)

type RegisterData struct {
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Password    string    `json:"password"`
	City        string    `json:"city"`
	Governorate string    `json:"governorate"`
	Address     string    `json:"address"`
	Birthdate   time.Time `json:"birthdate"`
}

func (s *Service) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var data RegisterData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		httputil.Error(nil, "invalid request body", http.StatusBadRequest).WriteJSON(w)
		return
	}
	defer r.Body.Close()

	user := models.User{
		Name:        data.Name,
		Email:       data.Email,
		Phone:       data.Phone,
		Role:        auth.RoleMember,
		City:        data.City,
		Governorate: data.Governorate,
		Address:     data.Address,
		Birthdate:   data.Birthdate,
	}
	if err := user.Validate(); err != nil {
		httputil.Error(nil, err.Error(), http.StatusBadRequest).WriteJSON(w)
		return
	}

	res := AuthResponse{}
	ctx := r.Context()
	err := s.db.InTx(ctx, func(tx pgx.Tx) error {
		if exists, err := models.UserExistsWithEmail(ctx, tx, user.Email); err != nil {
			return httputil.Error(err, "failed to get user by email", http.StatusInternalServerError)
		} else if exists {
			return httputil.Error(nil, "user with this email already exists", http.StatusConflict)
		}

		if exists, err := models.UserExistsWithPhone(ctx, tx, user.Phone); err != nil {
			return httputil.Error(err, "failed to get user by phone", http.StatusInternalServerError)
		} else if exists {
			return httputil.Error(nil, "user with this phone number already exists", http.StatusConflict)
		}

		userID, err := s.auth.GenerateID()
		if err != nil {
			return httputil.Error(err, "failed to generate id", http.StatusInternalServerError)
		}
		user.ID = userID

		hash, salt, err := s.auth.Hash(user.Password)
		if err != nil {
			return httputil.Error(err, "failed to hash password", http.StatusInternalServerError)
		}
		user.Password = hash + "$" + salt

		if err := models.UserInsert(ctx, tx, &user); err != nil {
			return httputil.Error(err, "failed to create user", http.StatusInternalServerError)
		}

		sessionID, err := s.auth.GenerateID()
		if err != nil {
			return httputil.Error(err, "failed to generate id", http.StatusInternalServerError)
		}

		refreshToken, err := s.auth.GenerateToken()
		if err != nil {
			return httputil.Error(err, "failed to generate token", http.StatusInternalServerError)
		}

		session := models.UserSession{
			ID:        sessionID,
			UserID:    userID,
			Token:     refreshToken,
			IPAddress: httputil.IP(r),
			UserAgent: httputil.UserAgent(r),
		}
		if err := models.UserSessionInsert(ctx, tx, &session); err != nil {
			return httputil.Error(err, "failed to create session", http.StatusInternalServerError)
		}

		client := httputil.Client(r)
		jwt, err := s.auth.GenerateJWT(user.ID, client, user.Role)
		if err != nil {
			return httputil.Error(err, "failed to generate jwt", http.StatusInternalServerError)
		}

		res.UserID = userID
		res.AccessToken = jwt
		res.RefreshToken = refreshToken
		res.ExpiresAt = session.ExpiresAt.Unix()

		return nil
	})
	if err != nil {
		if httpErr := httputil.CastError(err); httpErr != nil {
			s.logger.ErrorContext(ctx, "error occured in register handler", "error", httpErr.Wrap())
			httpErr.WriteJSON(w)
			return
		}
		s.logger.ErrorContext(ctx, "error occured in register handler", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&res)
}
