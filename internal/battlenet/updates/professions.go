package updates

import (
	"context"
	"fmt"
	"sync"

	"github.com/levifitzpatrick1/guild_tracker/generated/dbstore"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet"
)

// Updates the characters professions and recipes
// based on name and server. it passes a pointer to the
// recipe map to help reduce lookups.
func UpdateCharacterProfessions(ctx context.Context, name, server string, b *battlenet.Battlenet, q *dbstore.Queries, recipeCache *sync.Map) error {
	character, err := b.GetCharacterProfessions(name, server)
	if err != nil {
		return fmt.Errorf("Failed to get character profession info: %w", err)
	}

	for _, profession := range character.Primaries {
		for _, tier := range profession.Tiers {
			for _, recipe := range tier.KnownRecipes {
				recipeID := int64(recipe.ID)

				if _, exists := recipeCache.Load(recipeID); !exists {

					_, err := q.UpsertRecipe(ctx, dbstore.UpsertRecipeParams{
						ID:         recipeID,
						Name:       recipe.Name,
						Profession: profession.Profession.Name,
					})
					if err != nil {
						b.Logger.Error("Failed to upsert recipe %s: %v", recipe.Name, err)
						continue
					}

					err = UpdateRecipeMaterials(ctx, recipeID, b, q)
					if err != nil {
						b.Logger.Error("Failed to update materials for recipe %s: %v", recipe.Name, err)
					}

					recipeCache.Store(recipeID, true)
				}

				err = q.UpsertCharacterRecipe(ctx, dbstore.UpsertCharacterRecipeParams{
					CharacterID: int64(character.Character.ID),
					RecipeID:    recipeID,
				})
				if err != nil {
					b.Logger.Error("Failed to link recipe %s to character: %v", recipe.Name, err)
				}
			}
		}
	}

	return nil
}
