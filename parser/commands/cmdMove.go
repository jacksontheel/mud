package commands

import "example.com/mud/models"

var moveCommand = CommandDefinition{
	Aliases: []string{"move", "go", "walk"},
	Pattern: []models.Pattern{
		{
			Kind: "move",
			Tokens: []models.PatToken{
				models.Slot("direction"),
			},
			NoMatchMessage: "You can't get there.",
		},

		{
			Kind: "move",
			Tokens: []models.PatToken{
				models.Lit("move"),
				models.SlotRest("direction"),
			},
			NoMatchMessage: "You can't get there.",
		},
	},
}
