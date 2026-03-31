package scheduler

import (
	"os"
	"time"
	_ "time/tzdata"

	"github.com/levifitzpatrick1/guild_tracker/generated/dbstore"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet"
	"github.com/levifitzpatrick1/guild_tracker/internal/wrappers"
	"github.com/robfig/cron/v3"
)

// Scheduler to 
type Scheduler struct {
	Logger    *wrappers.Logger
	Queries   *dbstore.Queries
	Battlenet *battlenet.Battlenet
	Cron      *cron.Cron
	GuildName string
	Server    string
}

func Start(l *wrappers.Logger, q *dbstore.Queries, b *battlenet.Battlenet) error {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return err
	}

	s := &Scheduler{
		Logger:    l,
		Queries:   q,
		Battlenet: b,
		Cron:      cron.New(cron.WithLocation(loc)),
		GuildName: os.Getenv("GUILD_NAME"),
		Server:    os.Getenv("GUILD_SERVER"),
	}

	if _, err := s.Cron.AddFunc("0 8 * * *", s.dailyGuildSync); err != nil {
		return err
	}

	s.Cron.Start()
	l.Info("Scheduler running...")
	return nil
}
