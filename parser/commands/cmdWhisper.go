package commands

import "example.com/mud/models"

var WhisperPatterns = []models.Pattern{
	{Kind: "whisper", Tokens: []models.PatToken{
		models.Lit("whisper"),
		models.Slot("target"),
		models.SlotRest("message")}},
}

var WhisperAliases = map[string]string{
	"whisper": "whisper",
}
