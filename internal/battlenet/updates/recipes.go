package updates

import (
	"context"
	"fmt"

	"github.com/levifitzpatrick1/guild_tracker/generated/dbstore"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet"
)

func UpdateRecipeMaterials(ctx context.Context, id int64, b *battlenet.Battlenet, q *dbstore.Queries) error {
	recipeDetails, err := b.GetRecipeData(id)
	if err != nil {
		return fmt.Errorf("Failed to get recipe details: %w", err)
	}

	for _, reagent := range recipeDetails.Reagents {

		_, err := q.UpsertMaterial(ctx, dbstore.UpsertMaterialParams{
			ID:   int64(reagent.Reagent.ID),
			Name: reagent.Reagent.Name,
		})
		if err != nil {
			b.Logger.Error("Failed to upsert material %s: %v", reagent.Reagent.Name, err)
			continue
		}

		err = q.UpsertRecipeMaterial(ctx, dbstore.UpsertRecipeMaterialParams{
			RecipeID:   id,
			MaterialID: int64(reagent.Reagent.ID),
			Quantity:   int64(reagent.Quantity),
		})
		if err != nil {
			b.Logger.Error("Failed to link material %s to recipe: %v", reagent.Reagent.Name, err)
		}
	}

	return nil
}
