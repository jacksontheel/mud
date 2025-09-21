package actions

import (
	"fmt"
	"strings"

	"example.com/mud/world/entities"
	"example.com/mud/world/entities/components"
)

type Print struct {
	Text      string
	EventRole EventRole
}

var _ entities.Action = &Print{}

func (p *Print) Id() entities.ActionType {
	return entities.ActionPrint
}

func (p *Print) Execute(ev *entities.Event) error {
	if ev.Publisher == nil {
		return fmt.Errorf("publisher in event may not be nil for print action")
	}

	var recipient *entities.Entity
	switch p.EventRole {
	case EventRoleSource:
		recipient = ev.Source
	case EventRoleInstrument:
		recipient = ev.Instrument
	case EventRoleTarget:
		recipient = ev.Target
	default:
		return fmt.Errorf("invalid role '%s' for print action", p.EventRole.String())
	}

	message, err := formatText(p.Text, ev)
	if err != nil {
		return err
	}

	ev.Publisher.PublishTo(ev.Room, recipient, message)

	return nil
}

// eventually this could be its own package, handling text colors, etc
// for now it just does some simple replacements.
func formatText(s string, ev *entities.Event) (string, error) {
	if ev.Source != nil {
		identity, err := entities.RequireComponent[*components.Identity](ev.Source)
		if err != nil {
			return "", fmt.Errorf("could not format {source} for event: %w", err)
		}

		s = strings.ReplaceAll(s, "{source}", identity.Name)
	}

	if ev.Instrument != nil {
		identity, err := entities.RequireComponent[*components.Identity](ev.Instrument)
		if err != nil {
			return "", fmt.Errorf("could not format {instrument} for event: %w", err)
		}

		s = strings.ReplaceAll(s, "{instrument}", identity.Name)
	}

	if ev.Target != nil {
		identity, err := entities.RequireComponent[*components.Identity](ev.Target)
		if err != nil {
			return "", fmt.Errorf("could not format {target} for event: %w", err)
		}

		s = strings.ReplaceAll(s, "{target}", identity.Name)
	}

	return s, nil
}
