package responseStructs

type RecipeDetails struct {
	Links                 Links                  `json:"_links"`
	ID                    int                    `json:"id"`
	Name                  string                 `json:"name"`
	Media                 Media                  `json:"media"`
	Reagents              []ReagentWrapper       `json:"reagents"`
	ModifiedCraftingSlots []ModifiedCraftingSlot `json:"modified_crafting_slots"`
}

type ReagentWrapper struct {
	Reagent  Reagent `json:"reagent"`
	Quantity int     `json:"quantity"`
}

type Reagent struct {
	Key  Link   `json:"key"`
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type ModifiedCraftingSlot struct {
	SlotType     SlotType `json:"slot_type"`
	DisplayOrder int      `json:"display_order"`
}

type SlotType struct {
	Key  Link   `json:"key"`
	Name string `json:"name"`
	ID   int    `json:"id"`
}
