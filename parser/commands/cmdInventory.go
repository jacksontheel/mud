package commands

import "example.com/mud/models"

var InventoryPatterns = []models.Pattern{
	{Kind: "inventory", Tokens: []models.PatToken{
		models.Lit("inventory")}},
}

var InventoryAliases = map[string]string{
	"inventory": "inventory",
	"i":         "inventory",
}
