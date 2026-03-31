package battlenet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet/responseStructs"
)

var guildData string = "/data/wow/guild/%s/%s"

// Fetches the guild profile based on the guild name
// and server
func (b *Battlenet) GetGuildData(name, server string) (*responseStructs.GuildProfile, error) {
	name = ConvertToSlug(name)
	server = ConvertToSlug(server)
	locale := fmt.Sprintf("en_%s", strings.ToUpper(b.Region))
	namespace := fmt.Sprintf("profile-%s", b.Region)

	uri := fmt.Sprintf(b.Uri+guildData, server, name)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", b.Token)
	q := req.URL.Query()
	q.Set("namespace", namespace)
	q.Add("locale", locale)
	req.URL.RawQuery = q.Encode()

	ret, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer ret.Body.Close()

	var guild responseStructs.GuildProfile
	if err := json.NewDecoder(ret.Body).Decode(&guild); err != nil {
		return nil, err
	}

	return &guild, err
}
