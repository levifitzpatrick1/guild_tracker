package wrappers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Function for a discord slash command
type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

// Grouping for a slash command and the function
// that should be run with it
type Route struct {
	Command *discordgo.ApplicationCommand
	Handler CommandHandler
}

// Router handles the registration and slash command
// use
type Router struct {
	Logger         *Logger
	Routes         map[string]Route
	session        *discordgo.Session
	guildID        string
	registeredCmds []*discordgo.ApplicationCommand
}

// Creates a new Router with no commands
func NewRouter(logger *Logger) *Router {
	return &Router{
		Logger: logger,
		Routes: make(map[string]Route),
	}
}

// Registers a route with the router
func (r *Router) Handle(cmd *discordgo.ApplicationCommand, handler CommandHandler) {
	r.Routes[cmd.Name] = Route{Command: cmd, Handler: handler}
}

// Creates all the handled commands with discord and stores the session
// and guildID so that we can clean it up later. Must be called after the passed in
// session has opened (s.Open())
func (r *Router) Register(s *discordgo.Session, gid string) error {
	r.session = s
	r.guildID = gid
	for _, route := range r.Routes {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, gid, route.Command)
		if err != nil {
			return fmt.Errorf("cannot create `%v` command: %w", route.Command.Name, err)
		}
		r.registeredCmds = append(r.registeredCmds, cmd)
	}
	return nil
}

// Runs the slash command that was sent from discord. Should be passed into
// s.addHandler()
func (r *Router) OnInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	cmdName := i.ApplicationCommandData().Name

	if route, exists := r.Routes[cmdName]; exists {
		r.Logger.Info("Executing slash command: /%s by %s", cmdName, i.Member.User.Username)
		route.Handler(s, i)
	}
}

// Delete all registered slash commands from discord. Shoud be deferred after Register
// is called
func (r *Router) Close() {
	for _, cmd := range r.registeredCmds {
		if err := r.session.ApplicationCommandDelete(r.session.State.User.ID, r.guildID, cmd.ID); err != nil {
			r.Logger.Error("cannot delete '%v' command: %v", cmd.Name, err)
		}
	}
}
