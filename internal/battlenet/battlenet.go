package battlenet

import (
	"fmt"
	"os"
	"strings"

	"github.com/levifitzpatrick1/guild_tracker/internal/wrappers"
)

type Battlenet struct {
	Logger *wrappers.Logger
	Uri    string
	Region string
	Token  string
}

// Creates a new Battlenet struct using env variables for
// the region. Gets a fresh token.
func New(logger *wrappers.Logger) (*Battlenet, error) {
	region := os.Getenv("REGION")
	if region == "" {
		region = "us"
	}

	uri := fmt.Sprintf("https://%s.api.blizzard.com", region)
	if region == "cn" {
		uri = "https://gateway.battlenet.com.cn"
	}

	token, err := getToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get initial api token: %w", err)
	}
	token = "Bearer " + token

	logger.Info("Created battlenet client with uri: %s", uri)

	return &Battlenet{
		Logger: logger,
		Uri:    uri,
		Region: region,
		Token:  token,
	}, nil

}

// Gets a new token and replaces the old one.
func (b *Battlenet) RefreshToken() error {
	token, err := getToken()
	if err != nil {
		return fmt.Errorf("Failed to refresh token: %w", err)
	}
	b.Token = "Bearer " + token

	b.Logger.Info("Refreshed battlenet token")

	return nil
}

// Converts a standard server name or character
// into a slug for the api (all lowercase, '-' replaces ' ')
func ConvertToSlug(slug string) string {
	returnSlug := strings.ToLower(slug)
	returnSlug = strings.ReplaceAll(returnSlug, " ", "-")

	return returnSlug
}
