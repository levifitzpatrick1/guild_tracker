package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	LogFilePath    string
	DiscordToken   string
	DiscordGuildID string
	Region         string
	BnetClientID   string
	BnetSecret     string
	GuildName      string
	GuildServer    string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		LogFilePath:    getEnvOrDefault("LOG_FILE_PATH", "./logs"),
		DiscordToken:   os.Getenv("DISCORD_TOKEN"),
		DiscordGuildID: os.Getenv("DISCORD_GUILD_ID"),
		Region:         getEnvOrDefault("REGION", "us"),
		BnetClientID:   os.Getenv("BATTLENET_CLIENT_ID"),
		BnetSecret:     os.Getenv("BATTLENET_CLIENT_SECRET"),
		GuildName:      os.Getenv("GUILD_NAME"),
		GuildServer:    os.Getenv("GUILD_SERVER"),
	}

	if cfg.DiscordToken == "" || cfg.DiscordGuildID == "" {
		return nil, fmt.Errorf("missing required discord credentials")
	}

	if cfg.BnetClientID == "" || cfg.BnetSecret == "" {
		return nil, fmt.Errorf("missing required Battle.net credentials")
	}

	if cfg.GuildName == "" || cfg.GuildServer == "" {
		return nil, fmt.Errorf("missing required guild tracking information")
	}

	return cfg, nil
}

func getEnvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
