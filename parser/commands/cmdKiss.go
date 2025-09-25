package commands

import "example.com/mud/models"

var KissPatterns = []models.Pattern{
	{Kind: "kiss", Tokens: []models.PatToken{
		models.Lit("kiss"),
		models.SlotRest("target")}},
}

var KissAliases = map[string]string{
	"kiss":    "kiss",
	"smooch":  "kiss",
	"makeout": "kiss",
}
