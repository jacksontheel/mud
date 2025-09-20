package actions

import "example.com/mud/world/entities"

type Whisper struct {
	Target string `json:"target"`
	Text   string `json:"text"`
}

var _ entities.Action = &Whisper{}

func (a *Whisper) Id() entities.ActionType {
	return entities.ActionWhisper
}

func (a *Whisper) Execute(ev *entities.Event) (string, bool) {
	return a.Text, true
}
