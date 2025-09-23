package cmd

import "example.com/mud/models"

var SayPatterns = []models.Pattern{
	{Kind: "say", Tokens: []models.PatToken{
		models.Lit("say"),
		models.SlotRest("message", "message")}},
}

var SayAliases = map[string]string{
	"say": "say",
}
