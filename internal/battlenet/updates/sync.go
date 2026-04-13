package updates

import (
	"context"
	"fmt"
	"sync"

	"github.com/levifitzpatrick1/guild_tracker/generated/dbstore"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet/responseStructs"
	"github.com/levifitzpatrick1/guild_tracker/internal/wrappers"
)

// Updates all the characters in the guild.
// for those who are above 90, it updates their recipes.
// Recipes are cached per run of sync guild to help reduce
// the amount of api calls.
func SyncGuild(ctx context.Context, guildName, server string, b *battlenet.Battlenet, q *dbstore.Queries, l *wrappers.Logger) error {
	members, err := UpdateCharactersInGuild(ctx, guildName, server, b, q)
	if err != nil {
		return fmt.Errorf("failed to sync roster: %w", err)
	}

	var recipeCache sync.Map
	var wg sync.WaitGroup
	sem := make(chan struct{}, 20)
	var mu sync.Mutex

	for _, member := range members {
		if member.Character.Level < 79 {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(m responseStructs.Member) {
			defer wg.Done()
			defer func() { <-sem }()

			err := UpdateCharacterProfessions(ctx, m.Character.Name, m.Character.Realm.Slug, b, q, &recipeCache)
			if err != nil {
				if err.Error() == "Failed to get character profession info: profile private or inactive" {
					l.Info("Skipping %s: profile is private or inactive.", m.Character.Name)
				} else {
					l.Error("Failed to update professions for %s: %v", m.Character.Name, err)
				}
			} else {
				mu.Lock()
				mu.Unlock()
			}
		}(member)
	}

	wg.Wait()
	return nil
}
