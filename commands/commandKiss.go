package commands

var kissPatterns = []Pattern{
	{Kind: "kiss", Tokens: []PatToken{Lit("kiss"), SlotRest("target", "target")}},
}

var kissAliases = map[string]string{
	"kiss":    "kiss",
	"smooch":  "kiss",
	"makeout": "kiss",
}
