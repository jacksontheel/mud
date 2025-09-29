package commands

import "example.com/mud/models"

var inventoryCommand = &models.CommandDefinition{
	Name:    "inventory",
	Aliases: []string{"inventory", "inv"},
	Patterns: []models.CommandPattern{
		{
			Slots: []string{}, // no slots needed
			Tokens: []models.PatToken{
				models.Lit("inventory"),
			},
			NoMatchMessage: "You can't check your inventory right now.",
		},
	},
}
