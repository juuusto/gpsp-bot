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
	);`)
	return err
}

func InsertVideoReaction(db *sql.DB, r VideoReaction) error {
	_, err := db.Exec(`
	INSERT INTO video_reactions (message_uuid, user, emoji, created_at)
	VALUES (?, ?, ?, ?)
	`, r.MessageUUID, r.User, r.Emoji, r.CreatedAt)
	return err
}

// GetReactionsForMessage returns all reactions for a specific message
func GetReactionsForMessage(db *sql.DB, messageUUID string) ([]VideoReaction, error) {
	rows, err := db.Query(`
	SELECT id, message_uuid, user, emoji, created_at
	FROM video_reactions
	WHERE message_uuid = ?
	ORDER BY created_at DESC
	`, messageUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []VideoReaction
	for rows.Next() {
		var r VideoReaction
		err := rows.Scan(&r.ID, &r.MessageUUID, &r.User, &r.Emoji, &r.CreatedAt)
		if err != nil {
			return nil, err
		}
		reactions = append(reactions, r)
	}
	return reactions, nil
}

// GetReactionStats returns statistics about reactions
func GetReactionStats(db *sql.DB) (map[string]int, error) {
	rows, err := db.Query(`
	SELECT emoji, COUNT(*) as count
	FROM video_reactions
	GROUP BY emoji
	ORDER BY count DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var emoji string
		var count int
		err := rows.Scan(&emoji, &count)
		if err != nil {
			return nil, err
		}
		stats[emoji] = count
	}
	return stats, nil
}

// GetUserReactions returns all reactions by a specific user
func GetUserReactions(db *sql.DB, username string) ([]VideoReaction, error) {
	rows, err := db.Query(`
	SELECT id, message_uuid, user, emoji, created_at
	FROM video_reactions
	WHERE user = ?
	ORDER BY created_at DESC
	`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []VideoReaction
	for rows.Next() {
		var r VideoReaction
		err := rows.Scan(&r.ID, &r.MessageUUID, &r.User, &r.Emoji, &r.CreatedAt)
		if err != nil {
			return nil, err
		}
		reactions = append(reactions, r)
	}
	return reactions, nil
}
