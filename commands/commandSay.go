package commands

import "example.com/mud/models"

var sayPatterns = []models.Pattern{
	{Kind: "say", Tokens: []models.PatToken{
		models.Lit("say"),
		models.SlotRest("message", "message")}},
}

var sayAliases = map[string]string{
	"say": "say",
}
