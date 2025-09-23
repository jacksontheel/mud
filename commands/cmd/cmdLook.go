package cmd

import "example.com/mud/models"

var LookPatterns = []models.Pattern{
	{Kind: "look", Tokens: []models.PatToken{
		models.Lit("look")}},

	{Kind: "look", Tokens: []models.PatToken{
		models.Lit("look"),
		models.Slot("target", "target")}},
}

var LookAliases = map[string]string{
	"look":    "look",
	"examine": "look",
	"inspect": "look",
	"l":       "look",
}
