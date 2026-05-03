package chatservice

import (
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/phoenix-industries/phoenix-server/pkg/database/models"
	"github.com/phoenix-industries/phoenix-server/pkg/httputil"
)

type SendDirectData struct {
	ToID    string `json:"to_id"`
	Content string `json:"content"`
}

func (d *SendDirectData) Validate() error {
	if d.ToID == "" {
		return httputil.NewStatusError(nil, "to_id is required", http.StatusBadRequest)
	}
	if d.Content == "" {
		return httputil.NewStatusError(nil, "content is required", http.StatusBadRequest)
	}
	return nil
}

func (s *Service) HandleSendDirect(w http.ResponseWriter, r *http.Request) *httputil.Response {
	userID, err := httputil.GetUserID(s.auth, r)
	if err != nil {
		return httputil.ErrUnauthorized.Response()
	}

	var data SendDirectData
	if err := httputil.BodyJSON(w, r, &data); err != nil {
		return httputil.ErrInvalidBody.Response()
	}
	if err := data.Validate(); err != nil {
		return httputil.ResponseFromError(err)
	}

	var room *models.ChatRoom
	ctx := r.Context()
	err = s.db.InTx(ctx, func(tx pgx.Tx) error {
		room, err = models.ChatRoomGetCommon(ctx, tx, userID, data.ToID)
		if err != nil {
			return httputil.NewStatusError(err, "failed to get chat room", http.StatusInternalServerError)
		}
		if room != nil {
			return nil
		}
		id, err := s.auth.GenerateID()
		if err != nil {
			return httputil.NewStatusError(err, "failed to create chat room", http.StatusInternalServerError)
		}
		room = &models.ChatRoom{
			Model: models.Model{
				ID: id,
			},
			Type:  models.ChatRoomTypeDirect,
			Title: "",
		}
		if err = models.ChatRoomInsert(ctx, tx, room); err != nil {
			return httputil.NewStatusError(err, "failed to create chat room", http.StatusInternalServerError)
		}
		fromMemberID, err := s.auth.GenerateID()
		if err != nil {
			return httputil.NewStatusError(err, "failed to create chat room", http.StatusInternalServerError)
		}
		toMemberID, err := s.auth.GenerateID()
		if err != nil {
			return httputil.NewStatusError(err, "failed to create chat room", http.StatusInternalServerError)
		}
		fromMember := models.ChatRoomMember{
			Model: models.Model{
				ID: fromMemberID,
			},
			RoomID: room.ID,
			UserID: userID,
			Banned: false,
		}
		if err = models.ChatRoomMemberInsert(ctx, tx, &fromMember); err != nil {
			return httputil.NewStatusError(err, "failed to create chat room", http.StatusInternalServerError)
		}
		toMember := models.ChatRoomMember{
			Model: models.Model{
				ID: toMemberID,
			},
			RoomID: room.ID,
			UserID: data.ToID,
			Banned: false,
		}
		if err = models.ChatRoomMemberInsert(ctx, tx, &toMember); err != nil {
			return httputil.NewStatusError(err, "failed to create chat room", http.StatusInternalServerError)
		}
		return nil
	})
	if err != nil {
		return httputil.ResponseFromError(err)
	}

	s.server.AddRoom(room.ID, []string{userID, data.ToID})
	return httputil.NewResponseOK(http.StatusOK, room)
}
