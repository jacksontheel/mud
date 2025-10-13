package actions

import (
	"fmt"

	"example.com/mud/utils"
	"example.com/mud/world/entities"
)

type Print struct {
	Text      string
	EventRole entities.EventRole
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
	case entities.EventRoleSource:
		recipient = ev.Source
	case entities.EventRoleInstrument:
		recipient = ev.Instrument
	case entities.EventRoleTarget:
		recipient = ev.Target
	default:
		return fmt.Errorf("invalid role '%s' for print action", p.EventRole.String())
	}

	message, err := utils.FormatText(p.Text, fillFormatMap(ev))
	if err != nil {
		return err
	}

	ev.Publisher.PublishTo(ev.Room, recipient, message)

	return nil
}

func fillFormatMap(ev *entities.Event) map[string]string {
	out := make(map[string]string, 4)

	out[entities.EventRoleMessage.String()] = ev.Message

	if ev.Source != nil {
		out[entities.EventRoleSource.String()] = ev.Source.Name
	}
	if ev.Instrument != nil {
		out[entities.EventRoleInstrument.String()] = ev.Instrument.Name
	}
	if ev.Target != nil {
		out[entities.EventRoleTarget.String()] = ev.Target.Name
	}

	return out
}
