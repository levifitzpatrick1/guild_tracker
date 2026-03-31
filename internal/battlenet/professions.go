package battlenet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet/responseStructs"
)

var characterProfessions string = "/profile/wow/character/%s/%s/professions"

// Fetches the professions of a specific character given their name and
// server
func (b *Battlenet) GetCharacterProfessions(name, server string) (*responseStructs.CharacterProfessions, error) {
	name = url.PathEscape(ConvertToSlug(name))
	server = ConvertToSlug(server)
	locale := fmt.Sprintf("en_%s", strings.ToUpper(b.Region))
	namespace := fmt.Sprintf("profile-%s", b.Region)

	uri := fmt.Sprintf(b.Uri+characterProfessions, server, name)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", b.Token)
	q := req.URL.Query()
	q.Set("namespace", namespace)
	q.Add("locale", locale)
	req.URL.RawQuery = q.Encode()

	var ret *http.Response

	for i := range 5 {
		ret, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		if ret.StatusCode == 429 {
			ret.Body.Close()
			time.Sleep(time.Duration(1<<i) * time.Second)
			continue
		}
		break
	}

	defer ret.Body.Close()

	if ret.StatusCode != http.StatusOK {
		b.Logger.Info("uri: %s", req.URL.Path)
		return nil, fmt.Errorf("bnet API returned status: %d", ret.StatusCode)
	}

	var professions responseStructs.CharacterProfessions
	if err := json.NewDecoder(ret.Body).Decode(&professions); err != nil {
		return nil, err
	}

	return &professions, nil

}
