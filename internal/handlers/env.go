package handlers

import (
	"github.com/levifitzpatrick1/guild_tracker/generated/dbstore"
	"github.com/levifitzpatrick1/guild_tracker/internal/wrappers"
)

type Env struct {
	DB     *dbstore.Queries
	Logger *wrappers.Logger
}

func New(db *dbstore.Queries, logger *wrappers.Logger) *Env {
	return &Env{
		DB:     db,
		Logger: logger,
	}
}
