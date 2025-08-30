package db

import (
	"database/sql"
)

var globalDB *sql.DB

func SetGlobalDB(db *sql.DB) {
	globalDB = db
}

func GetGlobalDB() *sql.DB {
	return globalDB
}