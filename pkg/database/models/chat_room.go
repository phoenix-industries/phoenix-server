package models

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/phoenix-industries/phoenix-server/pkg/database"
)

type ChatRoomType string

const (
	ChatRoomTypeDirect ChatRoomType = "direct"
	ChatRoomTypeGroup  ChatRoomType = "group"
)

type ChatRoom struct {
	Model
	Title string       `db:"title" json:"title"`
	Type  ChatRoomType `db:"type" json:"type"`
}

func (m *ChatRoom) Validate() error {
	if m.ID == "" {
		return errors.New("id is required")
	}
	if m.Title == "" && m.Type != ChatRoomTypeDirect {
		return errors.New("title is required")
	}
	if m.Type != ChatRoomTypeDirect && m.Type != ChatRoomTypeGroup {
		return errors.New("type is required")
	}
	return nil
}

func ChatRoomInsert(ctx context.Context, db database.DB, room *ChatRoom) error {
	if err := room.Validate(); err != nil {
		return err
	}
	query := `
		INSERT INTO chat_rooms
		(id, title, type)
		VALUES
		($1, $2, $3)
	`
	_, err := db.Exec(ctx, query, room.ID, room.Title, room.Type)
	return err
}

func ChatRoomGetByID(ctx context.Context, db database.DB, id string) (*ChatRoom, error) {
	query := `
		SELECT *
		FROM chat_rooms
		WHERE id = $1 AND deleted_at IS NULL
	`
	var room ChatRoom
	if err := pgxscan.Get(ctx, db, &room, query, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

type ChatRoomLastMessage struct {
	Content   *string    `db:"lm_content" json:"content,omitempty"`
	UserID    *string    `db:"lm_user_id" json:"user_id,omitempty"`
	UserName  *string    `db:"lm_user_name" json:"user_name,omitempty"`
	CreatedAt *time.Time `db:"lm_created_at" json:"created_at,omitempty"`
}

func (m *ChatRoomLastMessage) MarshalJSON() ([]byte, error) {
	if m == nil || m.UserID == nil || m.Content == nil {
		return []byte("null"), nil
	}
	return json.Marshal(m)
}

type ChatRoomListData struct {
	ChatRoom
	LastMessage *ChatRoomLastMessage `db:"" json:"last_message,omitempty"`
}

func ChatRoomList(ctx context.Context, db database.DB, userID string) ([]*ChatRoomListData, error) {
	query := `
		SELECT
			r.id,
			r.type,
			r.created_at,
			CASE r.type
				WHEN 'direct' THEN other_u.name
				ELSE r.title
			END                  AS title,
	        last_u.id            AS lm_user_id,
			last_u.name          AS lm_user_name,
			last_msg.content     AS lm_content,
			last_msg.created_at  AS lm_created_at
		FROM chat_room_members m
		INNER JOIN chat_rooms r ON r.id = m.room_id
		LEFT JOIN chat_room_members other_m ON other_m.room_id = r.id
			AND other_m.user_id != $1
			AND other_m.deleted_at IS NULL
			AND r.type = 'direct'
		LEFT JOIN users other_u ON other_u.id = other_m.user_id
		LEFT JOIN LATERAL (
			SELECT content, user_id, created_at
			FROM chat_messages
			WHERE room_id = r.id
			ORDER BY created_at DESC
			LIMIT 1
		) last_msg ON true
		LEFT JOIN users last_u ON last_u.id = last_msg.user_id
		WHERE m.user_id      = $1
			AND m.banned     = false
			AND m.deleted_at IS NULL
			AND r.deleted_at IS NULL
		ORDER BY last_msg.created_at DESC NULLS LAST;
	`
	var rooms []*ChatRoomListData
	if err := pgxscan.Select(ctx, db, &rooms, query, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return rooms, nil
}

func ChatRoomGetCommon(ctx context.Context, db database.DB, fromID, toID string) (*ChatRoom, error) {
	query := `
	SELECT r.*
	FROM chat_rooms r
	INNER JOIN chat_room_members m1 ON m1.room_id = r.id AND m1.user_id = $1
	INNER JOIN chat_room_members m2 ON m2.room_id = r.id AND m2.user_id = $2
	WHERE r.type = 'direct'
		AND r.deleted_at IS NULL
		AND m1.deleted_at IS NULL
		AND m2.deleted_at IS NULL
	LIMIT 1;
	`
	var room ChatRoom
	if err := pgxscan.Get(ctx, db, &room, query, fromID, toID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}
