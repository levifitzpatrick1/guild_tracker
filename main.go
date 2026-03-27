package main

import (
	"embed"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet"
	"github.com/levifitzpatrick1/guild_tracker/internal/handlers"
	"github.com/levifitzpatrick1/guild_tracker/internal/scheduler"
	"github.com/levifitzpatrick1/guild_tracker/internal/wrappers"
	_ "modernc.org/sqlite"
)

// embed the database migration scripts so that we
// can just run them when the bot starts, rather than
// having to deal with goose outside of the app
//
//go:embed database/migrations/*.sql
var embedMigrations embed.FS

func main() {
	// Load our enviormental variables
	godotenv.Load()

	lp := os.Getenv("LOG_FILE_PATH")
	if lp == "" {
		lp = "./logs"
	}
	t := os.Getenv("DISCORD_TOKEN")
	if t == "" {
		log.Fatal("No discord token provided...")
	}
	gid := os.Getenv("DISCORD_GUILD_ID")
	if gid == "" {
		log.Fatal("No server ID provided...")
	}

	// init the wrapper for the logger, with the path
	// from .env. This is just a way to handle multiwritters,
	// and potentially gotify later on for fatals
	l, err := wrappers.NewLogger(lp)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// Create the battlenet instance. This keeps track of
	// the region that the bot is being used for and handles
	// the api token for the blizzard api
	bn, err := battlenet.New(l)
	if err != nil {
		log.Fatalf("failed to make bnet: %v", err)
	}

	// connect to our database, make sure that its up to date
	// and return the queries from sqlc
	q, db, err := wrappers.Connect(embedMigrations)
	if err != nil {
		l.Fatal("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// use the discord token to connect to discord,
	// and create the discordgo wrapper
	dg, err := discordgo.New("Bot " + t)
	if err != nil {
		l.Fatal("Error creating discord session: %v", err)
	}

	// create the enviorment for the handlers. this just
	// lets us create funcs without passing in the queries
	// and logger to each of them.
	env := handlers.New(q, l)

	// the router handles the commands, the structure acting
	// similar to the mux router for http. We can register
	// commands from the 'handlers/' folder
	r := wrappers.NewRouter(q, l)

	// adding the ping command as a quick test
	r.Handle(handlers.PingCommand, env.Ping)

	// since our router handles all the interaction creation
	// we can just pass the routers on interaction into the
	// discordgo wrapper handler
	dg.AddHandler(r.OnInteraction)

	if err := dg.Open(); err != nil {
		l.Fatal("Error opening discord connection: %v", err)
	}
	
	if err := scheduler.Start(l, q, bn); err != nil {
		l.Error("Failed to start daily schediler: %v", err)
	}

	// discord requires us to register any slash commands so that
	// they actually show up on the server when you type the '/'.
	// This loops through our commands in the router and makes
	// sure that they get added
	l.Info("registering slash commands...")
	regCMDs := make([]*discordgo.ApplicationCommand, len(r.Commands))
	for i, v := range r.Commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, gid, v)
		if err != nil {
			l.Fatal("cannot create '%v' command: %v", v.Name, err)
		}
		regCMDs[i] = cmd
	}

	// basic don't end the bot until we hit 'ctrl+C'
	l.Info("bot running...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// when we shut down the bot, we clean up the commands from the
	// server
	l.Info("shutting down...")
	for _, v := range regCMDs {
		err := dg.ApplicationCommandDelete(dg.State.User.ID, gid, v.ID)
		if err != nil {
			l.Error("cannot delete '%v' command: %v", v.Name, err)
		}
	}

	dg.Close()

}
