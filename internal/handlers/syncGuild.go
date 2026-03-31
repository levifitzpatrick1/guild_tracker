package handlers

import (
	"context"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet/updates"
)

var SyncCommand = &discordgo.ApplicationCommand{
	Name:        "Sync Guild",
	Description: "Force a manual sync of the guild roster and professions",
}

// Updates all the characters in the guild in the env variables
func (e *Env) SyncGuild(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		e.Logger.Error("Failed to defer interaction: %v", err)
		return
	}

	guildName := os.Getenv("GUILD_NAME")
	server := os.Getenv("GUILD_SERVER")

	ctx := context.Background()

	if err := updates.SyncGuild(ctx, guildName, server, e.Bnet, e.DB, e.Logger); err != nil {
		e.Logger.Error("Sync failed: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to sync guild...",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Sync complete!",
		},
	})

}
