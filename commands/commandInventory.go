package commands

var inventoryPatterns = []Pattern{
	{Kind: "inventory", Tokens: []PatToken{Lit("inventory")}},
}

var inventoryAliases = map[string]string{
	"inventory": "inventory",
	"i":         "inventory",
}
