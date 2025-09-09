package entities

import "fmt"

type Event struct {
	Kind   string
	Target *Entity
}

type When struct {
	Kind string `json:"kind"`
}

type Rule struct {
	When When `json:"when"`
}

type Eventful interface {
	OnEvent(ev *Event) (string, bool)
}

type CEventful struct {
	Rules []Rule `json:"rules"`
}

func (c *CEventful) OnEvent(ev *Event) (string, bool) {
	for _, r := range c.Rules {
		if r.When.Kind == ev.Kind {
			return fmt.Sprintf("match '%s'", ev.Kind), true
		}
	}
	return "no match", false
}
