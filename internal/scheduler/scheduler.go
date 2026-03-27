package scheduler

import (
	"context"
	"os"
	"time"

	"github.com/levifitzpatrick1/guild_tracker/generated/dbstore"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet/updates"
	"github.com/levifitzpatrick1/guild_tracker/internal/wrappers"
	"github.com/robfig/cron/v3"
)

func Start(l *wrappers.Logger, q *dbstore.Queries, b *battlenet.Battlenet) error {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return err
	}

	c := cron.New(cron.WithLocation(loc))

	guildName := os.Getenv("GUILD_NAME")
	server := os.Getenv("GUILD_SERVER")

	_, err = c.AddFunc("0 8 * * *", func() {
		l.Info("Starting daily 8:00 AM EST sync...")
		ctx := context.Background()

		if err := b.RefreshToken(); err != nil {
			l.Error("failed to refresh bnet token on daily sync: %v", err)
			return
		}

		members, err := updates.UpdateCharactersInGuild(ctx, guildName, server, b, q)
		if err != nil {
			l.Error("Daily roster sync failed: %v", err)
		} else {
			l.Info("Daily roster sync completed. Processing %d members", len(members))
		}

		for _, member := range members {
			if member.Character.Level < 90 {
				continue
			}

			err := updates.UpdateCharacterProfessions(ctx, member.Character.Name, member.Character.Realm.Name, b, q)
			if err != nil {
				l.Error("Failed to update professions for %s-%s: %v", member.Character.Name, member.Character.Realm.Name, err)
			}

			time.Sleep(50 * time.Millisecond)
		}
	})

	if err != nil {
		return err
	}

	c.Start()
	l.Info("Daily schedular running @ 8:00 AM EST.")

	return nil
}
