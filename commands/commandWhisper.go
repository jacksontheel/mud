package commands

import "example.com/mud/models"

var whisperPatterns = []models.Pattern{
	{Kind: "whisper", Tokens: []models.PatToken{
		models.Lit("whisper"),
		models.Slot("target", "target"),
		models.SlotRest("message", "message")}},
}

var whisperAliases = map[string]string{
	"whisper": "whisper",
}
