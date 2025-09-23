package cmd

import "example.com/mud/models"

var MovePatterns = []models.Pattern{
	{Kind: "move", Tokens: []models.PatToken{
		models.Slot("direction", "direction")}},

	{Kind: "move", Tokens: []models.PatToken{
		models.Lit("move"),
		models.SlotRest("direction", "direction")}},
}

var MoveAliases = map[string]string{
	"move": "move",
	"go":   "move",
	"walk": "move",
}

var MoveDirectionAliases = map[string]string{
	"n":     models.DirectionNorth,
	"north": models.DirectionNorth,

	"s":     models.DirectionSouth,
	"south": models.DirectionSouth,

	"e":    models.DirectionEast,
	"east": models.DirectionEast,

	"w":    models.DirectionWest,
	"west": models.DirectionWest,

	"u":  models.DirectionUp,
	"up": models.DirectionUp,

	"d":    models.DirectionDown,
	"down": models.DirectionDown,
}
