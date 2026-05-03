package models

import (
	"context"
	"errors"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/phoenix-industries/phoenix-server/pkg/database"
)

type ChatMessageType string

const (
	ChatMessageTypeText  ChatMessageType = "text"
	ChatMessageTypeFile  ChatMessageType = "file"
	ChatMessageTypeImage ChatMessageType = "image"
	ChatMessageTypeAudio ChatMessageType = "audio"
	ChatMessageTypeVideo ChatMessageType = "video"
)

type ChatMessage struct {
	Model
	RoomID  string          `db:"room_id" json:"room_id"`
	UserID  string          `db:"user_id" json:"user_id"`
	Type    ChatMessageType `db:"type" json:"type"`
	Content string          `db:"content" json:"content"`
}

func (m *ChatMessage) Validate() error {
	if m.ID == "" {
		return errors.New("id is required")
	}
	if m.RoomID == "" {
		return errors.New("room_id is required")
	}
	if m.UserID == "" {
		return errors.New("user_id is required")
	}
	switch m.Type {
	case ChatMessageTypeText, ChatMessageTypeFile, ChatMessageTypeImage, ChatMessageTypeAudio, ChatMessageTypeVideo:
		// ok
	default:
		return errors.New("type is required")
	}
	if m.Content == "" {
		return errors.New("content is required")
	}
	return nil
}

func ChatMessageInsert(ctx context.Context, db database.DB, message *ChatMessage) error {
	if err := message.Validate(); err != nil {
		return err
	}
	query := `
		INSERT INTO chat_messages
		(id, room_id, user_id, type, content)
		VALUES
		($1, $2, $3, $4)
	`
	_, err := db.Exec(ctx, query, message.ID, message.RoomID, message.UserID, message.Type, message.Content)
	return err
}

func ChatMessageList(ctx context.Context, db database.DB, roomID, userID string, before time.Time, limit int) ([]*ChatMessage, error) {
	if limit <= 0 || limit > 50 {
		limit = 50
	}
	query := `
		SELECT
			m.id,
			m.content,
			m.created_at,
			u.id    AS user_id,
			u.name  AS user_name
		FROM chat_messages m
		JOIN users u ON u.id = m.user_id
		WHERE m.room_id = $2
			AND m.created_at < $4
			AND EXISTS (
				SELECT 1
				FROM chat_room_members
				WHERE room_id     = $2
					AND user_id   = $3
					AND banned    = false
					AND deleted_at IS NULL
		  )
		ORDER BY m.created_at DESC
		LIMIT $1;
	`
	var messages []*ChatMessage
	if err := pgxscan.Select(ctx, db, &messages, query, limit, roomID, userID, before); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return messages, nil
}
