package commands

var takePatterns = []Pattern{
	{Kind: "take", Tokens: []PatToken{Lit("take"), Slot("target", "target")}},
}

var takeAliases = map[string]string{
	"take":    "take",
	"get":     "take",
	"grab":    "take",
	"pickup":  "take",
	"collect": "take",
}
