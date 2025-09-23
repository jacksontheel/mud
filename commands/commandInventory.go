package commands

import "example.com/mud/models"

var inventoryPatterns = []models.Pattern{
	{Kind: "inventory", Tokens: []models.PatToken{
		models.Lit("inventory")}},
}

var inventoryAliases = map[string]string{
	"inventory": "inventory",
	"i":         "inventory",
}
