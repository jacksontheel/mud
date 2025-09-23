package commands

var sayPatterns = []Pattern{
	{Kind: "say", Tokens: []PatToken{Lit("say"), SlotRest("message", "message")}},
}

var sayAliases = map[string]string{
	"say": "say",
}
