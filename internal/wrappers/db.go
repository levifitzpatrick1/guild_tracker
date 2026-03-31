package wrappers

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/levifitzpatrick1/guild_tracker/generated/dbstore"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

// Connect to the local sqlite database store. Use goose to update
// the project with any data migrations, then return the sqlc queries
// and the database pointer.
func Connect(e embed.FS) (*dbstore.Queries, *sql.DB, error) {
	db, err := sql.Open("sqlite", "./database/app.db?_busy_timeout=5000&_journal_mode=WAL")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}
	db.SetMaxOpenConns(1)

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	goose.SetBaseFS(e)
	if err := goose.SetDialect("sqlite"); err != nil {
		return nil, nil, fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db, "database/migrations"); err != nil {
		return nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return dbstore.New(db), db, nil
}
