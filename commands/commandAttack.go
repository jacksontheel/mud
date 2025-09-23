package commands

import "example.com/mud/models"

var attackPatterns = []models.Pattern{
	{Kind: "attack", Tokens: []models.PatToken{
		models.Lit("attack"),
		models.Slot("target", "target"),
		models.Lit("with"),
		models.Slot("instrument", "instrument")}},

	{Kind: "attack", Tokens: []models.PatToken{
		models.Lit("attack"),
		models.SlotRest("target", "target")}},
}

var attackAliases = map[string]string{
	"attack": "attack",
	"hit":    "attack",
	"kill":   "attack",
}
