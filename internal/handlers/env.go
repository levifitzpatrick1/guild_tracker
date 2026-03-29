package handlers

import (
	"github.com/levifitzpatrick1/guild_tracker/generated/dbstore"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet"
	"github.com/levifitzpatrick1/guild_tracker/internal/wrappers"
)

type Env struct {
	DB     *dbstore.Queries
	Logger *wrappers.Logger
	Bnet   *battlenet.Battlenet
}

func New(db *dbstore.Queries, logger *wrappers.Logger, b *battlenet.Battlenet) *Env {
	return &Env{
		DB:     db,
		Logger: logger,
		Bnet:   b,
	}
}
