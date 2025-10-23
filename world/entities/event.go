package entities

import (
	"fmt"
	"strconv"

	"example.com/mud/models"
	"example.com/mud/utils"
)

// allows us to use the bus without tightly coupling a
// specific publisher to our world model
type Publisher interface {
	Publish(room *Entity, text string, exclude []*Entity)
	PublishTo(room *Entity, recipient *Entity, text string)
}

type Event struct {
	Type         string
	Publisher    Publisher
	EntitiesById map[string]*Entity
	Room         *Entity
	Source       *Entity
	Instrument   *Entity
	Target       *Entity
	Message      string
}

type Rule struct {
	When []Condition
	Then []Action
}

func FormatEventMessage(message string, ev *Event) (string, error) {
	eventMap := make(map[string]string, 4)

	eventMap[EventRoleMessage.String()] = ev.Message

	if ev.Source != nil {
		role := EventRoleSource.String()
		eventMap[role] = ev.Source.Name

		for f, v := range ev.Source.Fields {
			switch v.K {
			case models.KindBool:
				eventMap[fmt.Sprintf("%s.%s", role, f)] = strconv.FormatBool(v.B)
			case models.KindInt:
				eventMap[fmt.Sprintf("%s.%s", role, f)] = strconv.FormatInt(int64(v.I), 10)
			case models.KindString:
				eventMap[fmt.Sprintf("%s.%s", role, f)] = v.S
			}
		}
	}

	if ev.Instrument != nil {
		role := EventRoleInstrument.String()
		eventMap[role] = ev.Instrument.Name

		for f, v := range ev.Instrument.Fields {
			switch v.K {
			case models.KindBool:
				eventMap[fmt.Sprintf("%s.%s", role, f)] = strconv.FormatBool(v.B)
			case models.KindInt:
				eventMap[fmt.Sprintf("%s.%s", role, f)] = strconv.FormatInt(int64(v.I), 10)
			case models.KindString:
				eventMap[fmt.Sprintf("%s.%s", role, f)] = v.S
			}
		}
	}

	if ev.Target != nil {
		role := EventRoleTarget.String()
		eventMap[role] = ev.Target.Name

		for f, v := range ev.Target.Fields {
			switch v.K {
			case models.KindBool:
				eventMap[fmt.Sprintf("%s.%s", role, f)] = strconv.FormatBool(v.B)
			case models.KindInt:
				eventMap[fmt.Sprintf("%s.%s", role, f)] = strconv.FormatInt(int64(v.I), 10)
			case models.KindString:
				eventMap[fmt.Sprintf("%s.%s", role, f)] = v.S
			}
		}
	}

	message, err := utils.FormatText(message, eventMap)
	if err != nil {
		return "", fmt.Errorf("format event message: %w", err)
	}

	return message, nil
}
