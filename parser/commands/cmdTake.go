package commands

import "example.com/mud/models"

var takeCommand = CommandDefinition{
	Aliases: []string{"take", "get", "grab", "pickup", "collect"},
	Pattern: []models.Pattern{
		{
			Kind: "take",
			Tokens: []models.PatToken{
				models.Lit("take"),
				models.Slot("target"),
			},
			NoMatchMessage: "You can't pick that up!",
		},

		// Take pattern for picking up multiple items --
		// {
		// 	Kind: "take",
		// 	Tokens: []models.PatToken{
		// 		models.Lit("take"),
		// 		models.SlotRest("target"),
		// 	},
		// 	NoMatchMessage: "You don't want to attack that.",
		// },

		// Take pattern for picking up all items in room --
	},
}
