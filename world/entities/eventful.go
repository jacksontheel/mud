package entities

import (
	"fmt"
	"strings"
)

type Event struct {
	Type   string
	Target *Entity
}

type When struct {
	Type string `json:"type"`
}

type Rule struct {
	When When
	Then []Action
}

type Eventful interface {
	OnEvent(ev *Event) (string, bool)
}

var _ Eventful = &CEventful{}

type CEventful struct {
	Rules []Rule
}

func (c *CEventful) OnEvent(ev *Event) (string, bool) {
	for _, r := range c.Rules {
		if r.When.Type == ev.Type {
			var b strings.Builder
			for _, a := range r.Then {
				if response, ok := performAction(a); ok {
					b.WriteString(fmt.Sprintf("%v\n", response))
				}
			}
			// only match on first match, return after
			return strings.TrimSuffix(b.String(), "\n"), true
		}
	}
	return "", false
}

func performAction(action Action) (string, bool) {
	if action.Id() == "say" {
		return action.(*ASay).Say(), true
	}

	return "", false
}
