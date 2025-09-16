package actions

import "example.com/mud/world/entities"

type Say struct {
	Text string `json:"text"`
}

var _ entities.Action = &Say{}

func (a *Say) Id() entities.ActionType {
	return entities.ActionSay
}

func (a *Say) Execute(ev *entities.Event) (string, bool) {
	return a.Text, true
}
