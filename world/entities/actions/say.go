package actions

import "example.com/mud/world/entities"

type Say struct {
	Text string `json:"text"`
}

var _ entities.Action = &Say{}

func (a *Say) Id() string {
	return "say"
}

func (a *Say) Execute(ev *entities.Event) (string, bool) {
	return a.Text, true
}
