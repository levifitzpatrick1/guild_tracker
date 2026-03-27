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

func New(l *wrappers.Logger) (*Battlenet, error) {
	r := os.Getenv("REGION")
	if r == "" {
		r = "us"
	}

	uri := fmt.Sprintf("https://%s.api.blizzard.com", r)
	if r == "cn" {
		uri = "https://gateway.battlenet.com.cn"
	}

	t, err := getToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get initial api token: %w", err)
	}
	t = "Bearer " + t

	l.Info("Created battlenet client with uri: %s", uri)

	return &Battlenet{
		Logger: l,
		Uri:    uri,
		Region: r,
		Token:  t,
	}, nil

}

func (b *Battlenet) RefreshToken() error {
	t, err := getToken()
	if err != nil {
		return fmt.Errorf("Failed to refresh token: %w", err)
	}
	b.Token = "Bearer " + t

	b.Logger.Info("Refreshed battlenet token")

	return nil
}

func ConvertToSlug(s string) string {
	r := strings.ToLower(s)
	r = strings.ReplaceAll(r, " ", "-")

	return r
}
