package actions

import (
	"fmt"

	"example.com/mud/world/entities"
)

type Print struct {
	Text      string
	EventRole entities.EventRole
}

var _ entities.Action = &Print{}

func (p *Print) Execute(ev *entities.Event) error {
	if ev.Publisher == nil {
		return fmt.Errorf("publisher in event may not be nil for print action")
	}

	var recipient *entities.Entity
	switch p.EventRole {
	case entities.EventRoleSource:
		recipient = ev.Source
	case entities.EventRoleInstrument:
		recipient = ev.Instrument
	case entities.EventRoleTarget:
		recipient = ev.Target
	default:
		return fmt.Errorf("invalid role '%s' for print action", p.EventRole.String())
	}

	message, err := entities.FormatEventMessage(p.Text, ev)
	if err != nil {
		return err
	}

	ev.Publisher.PublishTo(ev.Room, recipient, message)

	return nil
}
