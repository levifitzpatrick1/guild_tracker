package main

import (
	"embed"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/levifitzpatrick1/guild_tracker/internal/handlers"
	"github.com/levifitzpatrick1/guild_tracker/internal/wrappers"
	_ "modernc.org/sqlite"
)

//go:embed database/migrations/*.sql
var embedMigrations embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	lp := os.Getenv("LOG_FILE_PATH")
	if lp == "" {
		lp = "./logs"
	}
	l, err := wrappers.New(lp)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	q, db, err := wrappers.Connect(embedMigrations)
	if err != nil {
		l.Fatal("Failed to connect to database: %v", err)
	}
	defer db.Close()

	t := os.Getenv("DISCORD_TOKEN")
	gid := os.Getenv("DISCORD_GUILD_ID")

	dg, err := discordgo.New("Bot" + t)
	if err != nil {
		l.Fatal("Error creating discord session: %v", err)
	}

	env := handlers.New(q, l)
	r := wrappers.NewRouter(q, l)

	r.Handle(handlers.PingCommand, env.Ping)

	dg.AddHandler(r.OnInteraction)

	if err := dg.Open(); err != nil {
		l.Fatal("Error opening discord connection: %v", err)
	}

	l.Info("registering slash commands...")
	regCMDs := make([]*discordgo.ApplicationCommand, len(r.Commands))
	for i, v := range r.Commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, gid, v)
		if err != nil {
			l.Fatal("cannot create '%v' command: %v", v.Name, err)
		}
		regCMDs[i] = cmd
	}

	l.Info("bot running...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	l.Info("shutting down...")
	for _, v := range regCMDs {
		err := dg.ApplicationCommandDelete(dg.State.User.ID, gid, v.ID)
		if err != nil {
			l.Error("cannot delete '%v' command: %v", v.Name, err)
		}
	}

	dg.Close()

}
