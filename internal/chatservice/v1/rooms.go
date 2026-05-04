package chatservice

import (
	"net/http"

	"github.com/phoenix-industries/phoenix-server/pkg/database/models"
	"github.com/phoenix-industries/phoenix-server/pkg/httputil"
)

func (s *Service) HandleListRooms(w http.ResponseWriter, r *http.Request) *httputil.Response {
	userID, err := httputil.GetUserID(s.auth, r)
	if err != nil {
		return httputil.ErrUnauthorized.Response()
	}
	rooms, err := models.ChatRoomList(r.Context(), s.db.Pool(), userID)
	if err != nil {
		return httputil.NewStatusError(err, "failed to get rooms", http.StatusInternalServerError).Response()
	}
	return httputil.NewResponseOK(http.StatusOK, rooms)
}
