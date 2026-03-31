package battlenet

import (
	"encoding/json"
	"net/http"
)

// Struct to collect the API response from bnet
// for the token
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expiration  int    `json:"expires_in"`
}

// Fetches the bearer token as a string using the oauth
// credentials saved in the env vars.
func getToken(user, secret string) (string, error) {
	uri := "https://oauth.battle.net/token"
	req, err := http.NewRequest("POST", uri, nil)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(user, secret)

	q := req.URL.Query()
	q.Set("grant_type", "client_credentials")
	req.URL.RawQuery = q.Encode()

	ret, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer ret.Body.Close()

	var tr TokenResponse
	if err := json.NewDecoder(ret.Body).Decode(&tr); err != nil {
		return "", err
	}

	return tr.AccessToken, nil
}
