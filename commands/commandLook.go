package commands

var lookPatterns = []Pattern{
	{Kind: "look", Tokens: []PatToken{Lit("look")}},
	{Kind: "look", Tokens: []PatToken{Lit("look"), Slot("target", "target")}},
}

var lookAliases = map[string]string{
	"look":    "look",
	"examine": "look",
	"inspect": "look",
	"l":       "look",
}
