package commands

import "example.com/mud/models"

var lookPatterns = []models.Pattern{
	{Kind: "look", Tokens: []models.PatToken{
		models.Lit("look")}},

	{Kind: "look", Tokens: []models.PatToken{
		models.Lit("look"),
		models.Slot("target", "target")}},
}

var lookAliases = map[string]string{
	"look":    "look",
	"examine": "look",
	"inspect": "look",
	"l":       "look",
}
