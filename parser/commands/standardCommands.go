package commands

import "example.com/mud/models"

var inventoryCommand = models.CommandDefinition{
	Name:    "inventory",
	Aliases: []string{"inventory", "inv", "i"},
	Patterns: []models.CommandPattern{
		{
			Tokens: []models.PatToken{
				models.Lit("inventory"),
			},
			NoMatchMessage: "You can't check your inventory right now.",
		},
	},
}

var lookCommand = models.CommandDefinition{
	Name:    "look",
	Aliases: []string{"look", "l", "examine", "inspect"},
	Patterns: []models.CommandPattern{
		{
			Tokens: []models.PatToken{
				models.Lit("look"),
			},
			NoMatchMessage: "This shouldn't be possible (source: lookCommand)",
		},

		{
			Tokens: []models.PatToken{
				models.Lit("look"),
				models.SlotRest("target"),
			},
			NoMatchMessage: "There's nothing remarkable about that.",
		},
	},
}

var moveCommand = models.CommandDefinition{
	Name:    "move",
	Aliases: []string{"move", "go", "walk"},
	Patterns: []models.CommandPattern{
		{
			Tokens: []models.PatToken{
				models.Slot("direction"),
			},
			NoMatchMessage: "You can't get there.",
		},

		{
			Tokens: []models.PatToken{
				models.Lit("move"),
				models.SlotRest("direction"),
			},
			NoMatchMessage: "You can't get there.",
		},
	},
}

var sayCommand = models.CommandDefinition{
	Name:    "say",
	Aliases: []string{"say"},
	Patterns: []models.CommandPattern{
		{
			Tokens: []models.PatToken{
				models.Lit("say"),
				models.SlotRest("message"),
			},
			NoMatchMessage: "That's not allowed! (Source: sayCommand)",
		},
	},
}

var whisperCommand = models.CommandDefinition{
	Name:    "whisper",
	Aliases: []string{"whisper"},
	Patterns: []models.CommandPattern{
		{
			Tokens: []models.PatToken{
				models.Lit("whisper"),
				models.Slot("target"),
				models.SlotRest("message"),
			},
			NoMatchMessage: "It doesn't seem to hear you.",
		},
	},
}

var mapCommand = models.CommandDefinition{
	Name:    "map",
	Aliases: []string{"map", "m"},
	Patterns: []models.CommandPattern{
		{
			Tokens: []models.PatToken{
				models.Lit("map"),
			},
		},
	},
}
