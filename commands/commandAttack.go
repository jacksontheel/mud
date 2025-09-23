package commands

var attackPatterns = []Pattern{
	{Kind: "attack", Tokens: []PatToken{Lit("attack"), Slot("target", "target"), Lit("with"), Slot("instrument", "instrument")}},
	{Kind: "attack", Tokens: []PatToken{Lit("attack"), SlotRest("target", "target")}},
}

var attackAliases = map[string]string{
	"attack": "attack",
	"hit":    "attack",
	"kill":   "attack",
}
