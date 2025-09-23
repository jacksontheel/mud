package commands

var MovePatterns = []Pattern{
	{Kind: "move", Tokens: []PatToken{Slot("direction", "direction")}},
	{Kind: "move", Tokens: []PatToken{Lit("move"), SlotRest("direction", "direction")}},
}

var MoveAliases = map[string]string{
	"go":   "move",
	"walk": "move",
	"move": "move",
}
