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

	logPath := os.Getenv("LOG_FILE_PATH")
	if logPath == "" {
		logPath = "./logs"
	}
	discordToken := os.Getenv("DISCORD_TOKEN")
	if discordToken == "" {
		log.Fatal("No discord token provided...")
	}
	serverID := os.Getenv("DISCORD_GUILD_ID")
	if serverID == "" {
		log.Fatal("No server ID provided...")
	}

	// init the wrapper for the logger, with the path
	// from .env. This is just a way to handle multiwritters,
	// and potentially gotify later on for fatals
	logWrapper, err := wrappers.NewLogger(logPath)
	if err != nil {
		log.Fatal(err)
	}
	defer logWrapper.Close()

	// Create the battlenet instance. This keeps track of
	// the region that the bot is being used for and handles
	// the api token for the blizzard api
	battlenet, err := battlenet.New(logWrapper)
	if err != nil {
		log.Fatalf("failed to make bnet: %v", err)
	}

	// connect to our database, make sure that its up to date
	// and return the queries from sqlc
	q, db, err := wrappers.Connect(embedMigrations)
	if err != nil {
		logWrapper.Fatal("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// use the discord token to connect to discord,
	// and create the discordgo wrapper
	bot, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		logWrapper.Fatal("Error creating discord session: %v", err)
	}
	defer bot.Close()

	// create the enviorment for the handlers. this just
	// lets us create funcs without passing in the queries
	// and logger to each of them.
	env := handlers.New(q, logWrapper, battlenet)

	// the router handles the commands, the structure acting
	// similar to the mux router for http. We can register
	// commands from the 'handlers/' folder
	router := wrappers.NewRouter(logWrapper)
	defer router.Close()

	// add our command handlers and functions to the
	// router
	router.Handle(handlers.PingCommand, env.Ping)
	router.Handle(handlers.CraftCommand, env.Craft)
	router.Handle(handlers.SyncCommand, env.SyncGuild)

	// since our router handles all the interaction creation
	// we can just pass the routers on interaction into the
	// discordgo wrapper handler
	bot.AddHandler(router.OnInteraction)

	if err := bot.Open(); err != nil {
		logWrapper.Fatal("Error opening discord connection: %v", err)
	}

	if err := scheduler.Start(logWrapper, q, battlenet); err != nil {
		logWrapper.Error("Failed to start daily schediler: %v", err)
	}

	// discord requires us to register any slash commands so that
	// they actually show up on the server when you type the '/'.
	// This loops through our commands in the router and makes
	// sure that they get added
	logWrapper.Info("registering slash commands...")

	// register the router and all the commands to the
	// bot and server
	if err := router.Register(bot, serverID); err != nil {
		logWrapper.Error("Failed to register commands: %v", err)
	}

	// basic don't end the bot until we hit 'ctrl+C'
	logWrapper.Info("bot running...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// since everything is a defer, we don't need to worry about
	// shutting anything down
	logWrapper.Info("shutting down...")
}
