package commands

import "example.com/mud/models"

var inventoryCommand = CommandDefinition{
	Aliases: []string{"inventory", "i"},
	Pattern: []models.Pattern{
		{
			Kind: "inventory",
			Tokens: []models.PatToken{
				models.Lit("inventory"),
			},
			NoMatchMessage: "This shouldn't be possible",
		},
	},
}
