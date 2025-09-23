package commands

var whisperPatterns = []Pattern{
	{Kind: "whisper", Tokens: []PatToken{Lit("whisper"), Slot("target", "target"), SlotRest("message", "message")}},
}

var whisperAliases = map[string]string{
	"whisper": "whisper",
}
