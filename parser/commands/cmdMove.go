package commands

import "example.com/mud/models"

var MovePatterns = []models.Pattern{
	{Kind: "move", Tokens: []models.PatToken{
		models.Slot("direction")}},

	{Kind: "move", Tokens: []models.PatToken{
		models.Lit("move"),
		models.SlotRest("direction")}},
}

var MoveAliases = map[string]string{
	"move": "move",
	"go":   "move",
	"walk": "move",
}
