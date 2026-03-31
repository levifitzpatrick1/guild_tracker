package handlers

import (
	"github.com/bwmarrin/discordgo"
)

var PingCommand = &discordgo.ApplicationCommand{
	Name:        "ping",
	Description: "Check if the bot and database are connected",
}

// Basic command to check if the bot is running properly
func (e *Env) Ping(s *discordgo.Session, i *discordgo.InteractionCreate) {
	e.Logger.Info("Ping command executed by %s", i.Member.User.Username)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong! The router is working and the environment is injected.",
		},
	})
}
