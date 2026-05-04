package chatservice

import (
	"context"
	"errors"
	"net/http"

	"github.com/coder/websocket"
	"github.com/phoenix-industries/phoenix-server/pkg/httputil"
)

func (s *Service) HandleConnect(w http.ResponseWriter, r *http.Request) {
	userID, err := httputil.GetUserID(s.auth, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, "failed to accept websocket connection", http.StatusBadRequest)
		return
	}

	c := s.server.NewClient(conn, userID)
	err = s.server.Subscribe(c)
	if err == nil {
		return
	}
	if errors.Is(err, context.Canceled) {
		return
	}
	if status := websocket.CloseStatus(err); status == websocket.StatusNormalClosure || status == websocket.StatusGoingAway {
		return
	}
	s.logger.ErrorContext(r.Context(), "subscribe", "error", err)
}
