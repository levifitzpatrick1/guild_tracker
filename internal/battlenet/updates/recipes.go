package updates

import (
	"context"
	"fmt"

	"github.com/levifitzpatrick1/guild_tracker/generated/dbstore"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet"
)

// Given a specific recipe wow id, update the required materials and
// quantities. This shouldn't change often but is run as a quick check.
func UpdateRecipeMaterials(ctx context.Context, id int64, b *battlenet.Battlenet, q *dbstore.Queries) error {
	recipeDetails, err := b.GetRecipeData(id)
	if err != nil {
		return fmt.Errorf("Failed to get recipe details: %w", err)
	}

	for _, slot := range recipeDetails.ModifiedCraftingSlots {
		err = q.UpsertRecipeCraftingSlot(ctx, dbstore.UpsertRecipeCraftingSlotParams{
			RecipeID:     id,
			SlotName:     slot.SlotType.Name,
			DisplayOrder: int64(slot.DisplayOrder),
		})
		if err != nil {
			b.Logger.Error("Failed to upsert crafting slot for recipe %d: %v", id, err)
		}
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
			continue
		}
	}

	return nil
}
