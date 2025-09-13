package entities

import (
	"fmt"
	"strings"
)

type Event struct {
	Type       string
	Source     *Entity
	Instrument *Entity
	Target     *Entity
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

type Eventful struct {
	Rules []Rule
}

var _ Component = &Eventful{}

func (e *Eventful) Id() string {
	return "eventful"
}

func (c *Eventful) OnEvent(ev *Event) (string, bool) {
	for _, r := range c.Rules {
		if matchWhen(r.When, ev) {
			var b strings.Builder
			for _, a := range r.Then {
				if response, ok := performAction(a, ev); ok {
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
	return w.Type == ev.Type &&
		matchEntityToSelector(w.Source, ev.Source, ev.Target) &&
		matchEntityToSelector(w.Instrument, ev.Instrument, ev.Target)
}

func matchEntityToSelector(selector *EntitySelector, target, listener *Entity) bool {
	if selector == nil {
		return true
	}

	if target == nil {
		return false
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

func performAction(action Action, ev *Event) (string, bool) {
	switch action.Id() {
	case "say":
		return action.(*Say).Text, true
	case "removeItemFromInventory":
		removeItemAction := action.(*RemoveItemFromInventory)

		inventoryOwner := findEntityFromSelector(removeItemAction.InventoryOwner, ev)
		item := findEntityFromSelector(removeItemAction.Item, ev)

		removeItemAction.RemoveItemFromInventory(inventoryOwner, item)
	}

	return "", false
}

func findEntityFromSelector(selector EntitySelector, ev *Event) *Entity {
	switch selector.Type {
	case "source":
		return ev.Source
	case "instrument":
		return ev.Instrument
	case "target":
		return ev.Target
	}

	fmt.Println("This shouldn't happen")
	return nil
}
