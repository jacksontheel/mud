package commands

// import (
// 	"example.com/mud/models"
// )

// var attackCommand = CommandDefinition{
// 	Aliases: []string{"attack", "hit", "kill"},
// 	Pattern: []models.Pattern{
// 		{
// 			Kind: "attack",
// 			Tokens: []models.PatToken{
// 				models.Lit("attack"),
// 				models.Slot("target"),
// 				models.Lit("with"),
// 				models.Slot("instrument"),
// 			},
// 			NoMatchMessage: "You don't want to attack that with that.",
// 		},

// 		{
// 			Kind: "attack",
// 			Tokens: []models.PatToken{
// 				models.Lit("attack"),
// 				models.SlotRest("target"),
// 			},
// 			NoMatchMessage: "You don't want to attack that.",
// 		},
// 	},
// }
