package scheduler

import (
	"context"
	"os"
	"sync"
	"time"
	_ "time/tzdata"

	"github.com/levifitzpatrick1/guild_tracker/generated/dbstore"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet/responseStructs"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet/updates"
	"github.com/levifitzpatrick1/guild_tracker/internal/wrappers"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	Logger    *wrappers.Logger
	Queries   *dbstore.Queries
	Battlenet *battlenet.Battlenet
	Cron      *cron.Cron
}

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
			return
		} else {
			l.Info("Daily roster sync completed. Processing %d members", len(members))
		}

		var recipeCache sync.Map
		var wg sync.WaitGroup
		sem := make(chan struct{}, 5)
		var processed int
		var mu sync.Mutex

		for _, member := range members {
			if member.Character.Level < 90 {
				continue
			}

			wg.Add(1)
			sem <- struct{}{}

			go func(m responseStructs.Member) {
				defer wg.Done()
				defer func() { <-sem }()

				err := updates.UpdateCharacterProfessions(
					ctx,
					m.Character.Name,
					m.Character.Realm.Slug,
					int64(m.Character.ID),
					b,
					q,
					&recipeCache,
				)

				if err != nil {
					if err.Error() == "Failed to get character profession info: profile private or inactive" {
						l.Info("Skipping %s: profile is private or inactive.", m.Character.Name)
					} else {
						l.Error("Failed to update professions for %s-%s: %v", m.Character.Name, m.Character.Realm.Slug, err)
					}
				} else {
					mu.Lock()
					processed++
					mu.Unlock()
				}
			}(member)
		}

		wg.Wait()
		l.Info("Daily sync complete! Successfully updated %d high-level characters.", processed)
	})

	if err != nil {
		return err
	}

	c.Start()
	l.Info("Daily schedular running @ 8:00 AM EST.")

	return nil
}
