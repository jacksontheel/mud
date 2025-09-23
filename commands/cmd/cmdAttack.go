package cmd

import "example.com/mud/models"

var AttackPatterns = []models.Pattern{
	{Kind: "attack", Tokens: []models.PatToken{
		models.Lit("attack"),
		models.Slot("target", "target"),
		models.Lit("with"),
		models.Slot("instrument", "instrument")}},

	{Kind: "attack", Tokens: []models.PatToken{
		models.Lit("attack"),
		models.SlotRest("target", "target")}},
}

var AttackAliases = map[string]string{
	"attack": "attack",
	"hit":    "attack",
	"kill":   "attack",
}
