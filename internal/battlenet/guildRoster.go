package battlenet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet/responseStructs"
)

var guildRoster string = "/data/wow/guild/%s/%s/roster"

// Fetches a collection of all the members of a guid
// based off of guild name and server
func (b *Battlenet) GetGuildRoster(name, server string) (*responseStructs.GuildRoster, error) {
	name = ConvertToSlug(name)
	server = ConvertToSlug(server)
	locale := fmt.Sprintf("en_%s", strings.ToUpper(b.Region))
	namespace := fmt.Sprintf("profile-%s", b.Region)

	uri := fmt.Sprintf(b.Uri+guildRoster, server, name)
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

	var roster responseStructs.GuildRoster
	if err := json.NewDecoder(ret.Body).Decode(&roster); err != nil {
		return nil, err
	}

	return &roster, err
}
