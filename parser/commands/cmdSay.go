package commands

import "example.com/mud/models"

var sayCommand = CommandDefinition{
	Aliases: []string{"say"},
	Pattern: []models.Pattern{
		{
			Kind: "say",
			Tokens: []models.PatToken{
				models.Lit("say"),
				models.SlotRest("message"),
			},
			NoMatchMessage: "That's not allowed! (Source: sayCommand)",
		},
	},
}
