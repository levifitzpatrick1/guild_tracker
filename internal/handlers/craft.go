package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var CraftCommand = &discordgo.ApplicationCommand{
	Name:        "craft",
	Description: "Find who can craft an item and the materials required",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "item",
			Description: "The exact name of the item to craft",
			Required:    true,
		},
	},
}

func (e *Env) Craft(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()

	itemName := i.ApplicationCommandData().Options[0].StringValue()

	recipes, err := e.DB.SearchRecipes(ctx, itemName)
	if err != nil || len(recipes) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Could not find a recipe for '%s' in the database.", itemName),
			},
		})
		return
	}

	recipe := recipes[0]

	characters, err := e.DB.GetCharactersForRecipe(ctx, recipe.ID)
	if err != nil {
		e.Logger.Error("Failed to fetch characters for recipe %d: %v", recipe.ID, err)
	}

	crafterNames := []string{}
	for _, char := range characters {
		crafterNames = append(crafterNames, char.Name)
	}

	materials, err := e.DB.GetMaterialsForRecipe(ctx, recipe.ID)
	if err != nil {
		e.Logger.Error("Failed to fetch materials for recipe %d: %v", recipe.ID, err)
	}

	matStrings := []string{}
	for _, mat := range materials {
		matStrings = append(matStrings, fmt.Sprintf("- %dx %s", mat.Quantity, mat.Name))
	}

	crafterStr := "None currently known."
	if len(crafterNames) > 0 {
		crafterStr = strings.Join(crafterNames, ", ")
	}

	matStr := "No materials recorded."
	if len(matStrings) > 0 {
		matStr = strings.Join(matStrings, "\n")
	}

	responseContent := fmt.Sprintf("**Recipe:** %s (%s)\n\n**Crafters:** %s\n\n**Materials Required:**\n%s",
		recipe.Name, recipe.Profession, crafterStr, matStr)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseContent,
		},
	})
}
