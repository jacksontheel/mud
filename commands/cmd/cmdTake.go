package cmd

import "example.com/mud/models"

var TakePatterns = []models.Pattern{
	{Kind: "take", Tokens: []models.PatToken{
		models.Lit("take"),
		models.Slot("target", "target")}},
}

var TakeAliases = map[string]string{
	"take":    "take",
	"get":     "take",
	"grab":    "take",
	"pickup":  "take",
	"collect": "take",
}
