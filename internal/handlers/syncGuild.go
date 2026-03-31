package handlers

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet/updates"
)

var SyncCommand = &discordgo.ApplicationCommand{
	Name:        "sync-guild",
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

	guildName := e.Config.GuildName
	server := e.Config.GuildServer

	ctx := context.Background()

	if err := updates.SyncGuild(ctx, guildName, server, e.Bnet, e.DB, e.Logger); err != nil {
		e.Logger.Error("Sync failed: %v", err)
		_, _ = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Failed to sync guild...",
		})
		return
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "Sync complete!",
	})
	if err != nil {
		e.Logger.Error("Failed to send followup message: %v", err)
	}
}
