package handlers

import (
	"github.com/levifitzpatrick1/guild_tracker/generated/dbstore"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet"
	"github.com/levifitzpatrick1/guild_tracker/internal/config"
	"github.com/levifitzpatrick1/guild_tracker/internal/wrappers"
)

// Wrapper for required vars for the handler commands
type Env struct {
	DB     *dbstore.Queries
	Logger *wrappers.Logger
	Bnet   *battlenet.Battlenet
	Config *config.Config
}

// Returns a new Env
func New(db *dbstore.Queries, logger *wrappers.Logger, b *battlenet.Battlenet, cfg *config.Config) *Env {
	return &Env{
		DB:     db,
		Logger: logger,
		Bnet:   b,
		Config: cfg,
	}
}
