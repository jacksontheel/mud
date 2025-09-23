package commands

var TakePatterns = []Pattern{
	{Kind: "take", Tokens: []PatToken{Lit("take"), Slot("target", "target")}},
}

var TakeAliases = map[string]string{
	"take":    "take",
	"get":     "take",
	"grab":    "take",
	"pickup":  "take",
	"collect": "take",
}
