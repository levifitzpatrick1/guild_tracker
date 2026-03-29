package handlers

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet/responseStructs"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet/updates"
)

var UpdateCommand = &discordgo.ApplicationCommand{
	Name:        "update",
	Description: "Force a manual sync of the guild roster and professions",
}

func (e *Env) Update(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		e.Logger.Error("Failed to defer interaction: %v", err)
		return
	}

	guildName := os.Getenv("GUILD_NAME")
	server := os.Getenv("GUILD_SERVER")

	go func() {
		ctx := context.Background()

		if err := e.Bnet.RefreshToken(); err != nil {
			e.Logger.Error("failed to refresh bnet token on manual sync: %v", err)
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: new("Failed to authenticate with Battle.net."),
			})
			return
		}

		members, err := updates.UpdateCharactersInGuild(ctx, guildName, server, e.Bnet, e.DB)
		if err != nil {
			e.Logger.Error("Manual roster sync failed: %v", err)
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: new("Failed to fetch guild roster."),
			})
			return
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

				err := updates.UpdateCharacterProfessions(ctx, m.Character.Name, m.Character.Realm.Slug, int64(m.Character.ID), e.Bnet, e.DB, &recipeCache)

				if err != nil {
					if err.Error() == "Failed to get character profession info: profile private or inactive" {
						e.Logger.Info("Skipping %s: profile is private or inactive.", m.Character.Name)
					} else {
						e.Logger.Error("Failed to update professions for %s: %v", m.Character.Name, err)
					}
				} else {
					mu.Lock()
					processed++
					mu.Unlock()
				}
			}(member)
		}

		wg.Wait()

		msg := fmt.Sprintf("Sync complete! Updated %d high-level characters from **%s**.", processed, guildName)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
	}()
}
