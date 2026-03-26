package wrappers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/levifitzpatrick1/guild_tracker/generated/dbstore"
)

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

type Router struct {
	DB       *dbstore.Queries
	Logger   *Logger
	Commands []*discordgo.ApplicationCommand
	Handlers map[string]CommandHandler
}

func NewRouter(db *dbstore.Queries, logger *Logger) *Router {
	return &Router{
		DB:       db,
		Logger:   logger,
		Commands: []*discordgo.ApplicationCommand{},
		Handlers: make(map[string]CommandHandler),
	}
}

func (r *Router) Handle(cmd *discordgo.ApplicationCommand, handler CommandHandler) {
	r.Commands = append(r.Commands, cmd)
	r.Handlers[cmd.Name] = handler
}

func (r *Router) OnInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	cmdName := i.ApplicationCommandData().Name

	if handler, exists := r.Handlers[cmdName]; exists {
		r.Logger.Info("Executing slash command: /%s by %s", cmdName, i.Member.User.Username)
		handler(s, i)
	}
}
