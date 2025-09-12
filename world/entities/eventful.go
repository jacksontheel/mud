package entities

import (
	"fmt"
	"strings"
)

type Event struct {
	Type   string
	Source *Entity
	Target *Entity
}

type EntitySelector struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type When struct {
	Type       string          `json:"type"`
	Source     *EntitySelector `json:"source"`
	Instrument *EntitySelector `json:"instrument"`
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
		if matchWhen(r.When, ev) {
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

func matchWhen(w When, ev *Event) bool {
	return w.Type == ev.Type && matchEntityToSelector(w.Source, ev.Source, ev.Target)
}

func matchEntityToSelector(selector *EntitySelector, target, listener *Entity) bool {
	if selector == nil {
		return true
	}

	switch selector.Type {
	case "self":
		return target == listener
	case "tag":
		for _, t := range target.getTags() {
			if selector.Value == t {
				return true
			}
		}
	default:
		return false
	}

	return false
}

func performAction(action Action) (string, bool) {
	if action.Id() == "say" {
		return action.(*ASay).Say(), true
	}

	return "", false
}
