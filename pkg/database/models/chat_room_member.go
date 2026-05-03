package models

import (
	"context"
	"errors"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/phoenix-industries/phoenix-server/pkg/database"
)

type ChatRoomMember struct {
	Model
	RoomID   string     `db:"room_id" json:"room_id"`
	UserID   string     `db:"user_id" json:"user_id"`
	Banned   bool       `db:"banned" json:"banned"`
	JoinedAt *time.Time `db:"joined_at" json:"joined_at"`
}

func (m *ChatRoomMember) Validate() error {
	if m.ID == "" {
		return errors.New("id is required")
	}
	if m.RoomID == "" {
		return errors.New("room_id is required")
	}
	if m.UserID == "" {
		return errors.New("user_id is required")
	}
	return nil
}

func ChatRoomMemberInsert(ctx context.Context, db database.DB, member *ChatRoomMember) error {
	if err := member.Validate(); err != nil {
		return err
	}
	query := `
		INSERT INTO chat_room_members
		(id, room_id, user_id, banned, joined_at)
		VALUES
		($1, $2, $3, FALSE, CURRENT_TIMESTAMP)
	`
	_, err := db.Exec(ctx, query, member.ID, member.RoomID, member.UserID)
	return err
}

func ChatRoomMemberGetByID(ctx context.Context, db database.DB, id string) (*ChatRoomMember, error) {
	query := `
		SELECT *
		FROM chat_room_members
		WHERE id = $1 AND deleted_at IS NULL
	`
	var member ChatRoomMember
	if err := pgxscan.Get(ctx, db, &member, query, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

func ChatRoomMemberBanByID(ctx context.Context, db database.DB, roomID string, memberID string) error {
	query := `
		UPDATE chat_room_members
		SET banned = TRUE
		WHERE room_id = $1 AND id = $2 AND deleted_at IS NULL
	`
	_, err := db.Exec(ctx, query, roomID, memberID)
	return err
}

func ChatRoomMemberBanByUserID(ctx context.Context, db database.DB, roomID string, userID string) error {
	query := `
		UPDATE chat_room_members
		SET banned = TRUE
		WHERE room_id = $1 AND user_id = $2 AND deleted_at IS NULL
	`
	_, err := db.Exec(ctx, query, roomID, userID)
	return err
}

func ChatRoomMemberList(ctx context.Context, db database.DB, roomID string) ([]*ChatRoomMember, error) {
	query := `
		SELECT *
		FROM chat_room_members
		WHERE room_id = $1 AND deleted_at IS NULL
	`
	var members []*ChatRoomMember
	if err := pgxscan.Select(ctx, db, &members, query, roomID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return members, nil
}

func ChatRoomMemberByUserID(ctx context.Context, db database.DB, userID string) ([]*ChatRoom, error) {
	query := `
		SELECT *
		FROM chat_rooms
		WHERE id IN (
			SELECT room_id
			FROM chat_room_members
			WHERE user_id = $1 AND deleted_at IS NULL
		) AND deleted_at IS NULL
	`
	var rooms []*ChatRoom
	if err := pgxscan.Select(ctx, db, &rooms, query, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return rooms, nil
}
