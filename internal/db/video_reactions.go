package db

import (
	"database/sql"
)

type VideoReaction struct {
	ID         int64  // autoincrement
	MessageUUID string // message uuid (video message)
	User       string // user who reacted
	Emoji      string // emoji used
	CreatedAt  int64  // unix timestamp
}

func InitReactionsTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS video_reactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		message_uuid TEXT,
		user TEXT,
		emoji TEXT,
		created_at INTEGER
	);
	`)
	return err
}

func InsertVideoReaction(db *sql.DB, r VideoReaction) error {
	_, err := db.Exec(`
	INSERT INTO video_reactions (message_uuid, user, emoji, created_at)
	VALUES (?, ?, ?, ?)
	`, r.MessageUUID, r.User, r.Emoji, r.CreatedAt)
	return err
}
