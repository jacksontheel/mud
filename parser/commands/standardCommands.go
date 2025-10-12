package commands

import "example.com/mud/models"

var helpCommand = models.CommandDefinition{
	Name:    "help",
	Aliases: []string{"help", "h"},
	Patterns: []models.CommandPattern{
		{
			Tokens: []models.PatToken{
				models.Lit("help"),
			},
			HelpMessage: "Get more details about all commands (you're doing great!).",
		},
		{
			Tokens: []models.PatToken{
				models.Lit("help"),
				models.Slot("command"),
			},
			HelpMessage: "Get details about a specific command type.",
		},
	},
}

var inventoryCommand = models.CommandDefinition{
	Name:    "inventory",
	Aliases: []string{"inventory", "inv", "i"},
	Patterns: []models.CommandPattern{
		{
			Tokens: []models.PatToken{
				models.Lit("inventory"),
			},
			HelpMessage:    "Check items within your inventory.",
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
			HelpMessage:    "Get details about the room you're inside of.",
			NoMatchMessage: "This shouldn't be possible (source: lookCommand)",
		},

		{
			Tokens: []models.PatToken{
				models.Lit("look"),
				models.SlotRest("target"),
			},
			HelpMessage:    "Get details about a specific item in the room you're inside of.",
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
			HelpMessage:    "Move to another room.",
			NoMatchMessage: "You can't get there.",
		},

		{
			Tokens: []models.PatToken{
				models.Lit("move"),
				models.SlotRest("direction"),
			},
			HelpMessage:    "Move to another room.",
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
			HelpMessage:    "Say something out loud to other player's in your room.",
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
			HelpMessage:    "Say something to a specific player in your room.",
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
			HelpMessage: "Check out a map of the surrounding rooms.",
		},
	},
}
