package commands

import "example.com/mud/models"

var whisperCommand = CommandDefinition{
	Aliases: []string{"whisper"},
	Pattern: []models.Pattern{
		{
			Kind: "whisper",
			Tokens: []models.PatToken{
				models.Lit("whisper"),
				models.Slot("target"),
				models.SlotRest("message"),
			},
			NoMatchMessage: "It doesn't seem to hear you.",
		},
	},
}
