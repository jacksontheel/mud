package actions

import (
	"fmt"
	"strconv"

	"example.com/mud/models"
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
		role := entities.EventRoleSource.String()
		out[role] = ev.Source.Name

		for f, v := range ev.Source.Fields {
			switch v.K {
			case models.KindBool:
				out[fmt.Sprintf("%s.%s", role, f)] = strconv.FormatBool(v.B)
			case models.KindInt:
				out[fmt.Sprintf("%s.%s", role, f)] = strconv.FormatInt(int64(v.I), 10)
			case models.KindString:
				out[fmt.Sprintf("%s.%s", role, f)] = v.S
			}
		}
	}

	if ev.Instrument != nil {
		role := entities.EventRoleInstrument.String()
		out[role] = ev.Instrument.Name

		for f, v := range ev.Instrument.Fields {
			switch v.K {
			case models.KindBool:
				out[fmt.Sprintf("%s.%s", role, f)] = strconv.FormatBool(v.B)
			case models.KindInt:
				out[fmt.Sprintf("%s.%s", role, f)] = strconv.FormatInt(int64(v.I), 10)
			case models.KindString:
				out[fmt.Sprintf("%s.%s", role, f)] = v.S
			}
		}
	}

	if ev.Target != nil {
		role := entities.EventRoleTarget.String()
		out[role] = ev.Target.Name

		for f, v := range ev.Target.Fields {
			switch v.K {
			case models.KindBool:
				out[fmt.Sprintf("%s.%s", role, f)] = strconv.FormatBool(v.B)
			case models.KindInt:
				out[fmt.Sprintf("%s.%s", role, f)] = strconv.FormatInt(int64(v.I), 10)
			case models.KindString:
				out[fmt.Sprintf("%s.%s", role, f)] = v.S
			}
		}
	}

	return out
}
