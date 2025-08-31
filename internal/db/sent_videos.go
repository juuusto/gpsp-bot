package db

import (
	"database/sql"
	_ "github.com/marcboeker/go-duckdb"
	"log"
)

type SentVideo struct {
	ID           int64  // autoincrement
	User         string // user who sent
	VideoHash    string // hash of the video file
	MessageUUID  string // new message uuid
	Service      string // telegram/discord
	CreatedAt    int64  // unix timestamp
}

func InitDB(path string) *sql.DB {
	db, err := sql.Open("duckdb", path)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS sent_videos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user TEXT,
		video_hash TEXT,
		message_uuid TEXT,
		service TEXT,
		created_at INTEGER
	);
	`)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	// Initialize reactions table
	err = InitReactionsTable(db)
	if err != nil {
		log.Fatalf("failed to create reactions table: %v", err)
	}

	return db
}

func InsertSentVideo(db *sql.DB, v SentVideo) error {
	_, err := db.Exec(`
	INSERT INTO sent_videos (user, video_hash, message_uuid, service, created_at)
	VALUES (?, ?, ?, ?, ?)
	`, v.User, v.VideoHash, v.MessageUUID, v.Service, v.CreatedAt)
	return err
}
