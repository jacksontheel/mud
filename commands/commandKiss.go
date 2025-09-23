package commands

import "example.com/mud/models"

var kissPatterns = []models.Pattern{
	{Kind: "kiss", Tokens: []models.PatToken{
		models.Lit("kiss"),
		models.SlotRest("target", "target")}},
}

var kissAliases = map[string]string{
	"kiss":    "kiss",
	"smooch":  "kiss",
	"makeout": "kiss",
}
