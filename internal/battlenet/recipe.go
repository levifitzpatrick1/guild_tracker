package battlenet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet/responseStructs"
)

var recipeData string = "/data/wow/recipe/%d"

// Fetches the information on a specifc recipe based
// off of its wow id.
func (b *Battlenet) GetRecipeData(id int64) (*responseStructs.RecipeDetails, error) {
	locale := fmt.Sprintf("en_%s", strings.ToUpper(b.Region))
	namespace := fmt.Sprintf("static-%s", b.Region)

	uri := fmt.Sprintf(b.Uri+recipeData, id)
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

	if ret.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bnet api returned status: %d for recipe %d", ret.StatusCode, id)
	}

	var recipeDetails responseStructs.RecipeDetails
	if err := json.NewDecoder(ret.Body).Decode(&recipeDetails); err != nil {
		return nil, err
	}

	return &recipeDetails, err
}
