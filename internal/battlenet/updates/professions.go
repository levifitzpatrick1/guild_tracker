package updates

import (
	"context"
	"fmt"

	"github.com/levifitzpatrick1/guild_tracker/generated/dbstore"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet"
)

func UpdateCharacterProfessions(ctx context.Context, name, server string, b *battlenet.Battlenet, q *dbstore.Queries) error {
	character, err := b.GetCharacterProfessions(name, server)
	if err != nil {
		return fmt.Errorf("Failed to get character profession info: %w", err)
	}

	for _, profession := range character.Primaries {
		for _, tier := range profession.Tiers {
			for _, recipe := range tier.KnownRecipes {

				_, err := q.UpsertRecipe(ctx, dbstore.UpsertRecipeParams{
					ID:         int64(recipe.ID),
					Name:       recipe.Name,
					Profession: profession.Profession.Name,
				})
				if err != nil {
					b.Logger.Error("Failed to upsert recipe %s: %v", recipe.Name, err)
					continue
				}

				err = q.UpsertCharacterRecipe(ctx, dbstore.UpsertCharacterRecipeParams{
					CharacterID: int64(character.Character.ID),
					RecipeID:    int64(recipe.ID),
				})
				if err != nil {
					b.Logger.Error("Failed to link recipe %s to character: %v", recipe.Name, err)
				}
			}
		}
	}

	return nil
}
