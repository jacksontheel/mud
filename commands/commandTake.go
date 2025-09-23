package commands

import "example.com/mud/models"

var takePatterns = []models.Pattern{
	{Kind: "take", Tokens: []models.PatToken{
		models.Lit("take"),
		models.Slot("target", "target")}},
}

var takeAliases = map[string]string{
	"take":    "take",
	"get":     "take",
	"grab":    "take",
	"pickup":  "take",
	"collect": "take",
}
