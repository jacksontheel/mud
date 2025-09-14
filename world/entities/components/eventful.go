package components

import (
	"strings"

	"example.com/mud/world/entities"
)

type Eventful struct {
	Rules []entities.Rule
}

var _ entities.Component = &Eventful{}

func (e *Eventful) Copy() entities.Component {
	return &Eventful{
		Rules: e.Rules,
	}
}

func (c *Eventful) OnEvent(ev *entities.Event) (string, bool) {
	for _, r := range c.Rules {
		if matchWhen(r.When, ev) {
			var b strings.Builder
			for _, a := range r.Then {
				if resp, ok := a.Execute(ev); ok && resp != "" {
					if b.Len() > 0 {
						b.WriteByte('\n')
					}
					b.WriteString(resp)
				}
			}
			// only match on first match, return after
			return b.String(), true
		}
	}
	return "", false
}

func matchWhen(w entities.When, ev *entities.Event) bool {
	return w.Type == ev.Type &&
		matchEntityToSelector(w.Source, ev.Source, ev.Target) &&
		matchEntityToSelector(w.Instrument, ev.Instrument, ev.Target)
}

func matchEntityToSelector(selector *entities.EntitySelector, target, listener *entities.Entity) bool {
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
		for _, t := range GetTags(target) {
			if selector.Value == t {
				return true
			}
		}
	default:
		return false
	}

	return false
}
