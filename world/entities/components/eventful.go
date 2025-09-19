package components

import (
	"strings"

	"example.com/mud/world/entities"
)

type Eventful struct {
	Rules []*entities.Rule
}

var _ entities.Component = &Eventful{}

func (e *Eventful) Id() entities.ComponentType {
	return entities.ComponentEventful
}

func (e *Eventful) Copy() entities.Component {
	return &Eventful{
		Rules: e.Rules,
	}
}

func (c *Eventful) OnEvent(ev *entities.Event) (string, error) {
	for _, r := range c.Rules {
		match, err := matchWhen(r.When, ev)
		if err != nil {
			return "", err
		}

		if match {
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
			return b.String(), nil
		}
	}
	return "", nil
}

func (c *Eventful) AddRule(rule *entities.Rule) {
	c.Rules = append(c.Rules, rule)
}

func matchWhen(w *entities.When, ev *entities.Event) (bool, error) {

	sourceMatch, err := matchEntityToSelector(w.Source, ev.Source, ev.Target)
	if err != nil {
		return false, err
	}

	instrumentMatch, err := matchEntityToSelector(w.Instrument, ev.Instrument, ev.Target)
	if err != nil {
		return false, err
	}

	return w.Type == ev.Type && sourceMatch && instrumentMatch, nil
}

func matchEntityToSelector(selector *entities.EntitySelector, target, listener *entities.Entity) (bool, error) {
	// if there is no selector for target
	// default true
	if selector == nil {
		return true, nil
	}

	// if there is a selector for target
	// but target entity is nil
	// e.g. source selector with no event source
	// default false
	if target == nil {
		return false, nil
	}

	switch selector.Type {
	case "self":
		return target == listener, nil
	case "tag":
		identity, err := entities.RequireComponent[*Identity](target)
		if err != nil {
			return false, err
		}

		for _, t := range identity.Tags {
			if selector.Value == t {
				return true, nil
			}
		}
	default:
		return false, nil
	}

	return false, nil
}
