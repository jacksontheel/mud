package actions

import (
	"fmt"

	"example.com/mud/world/entities"
)

type Copy struct {
	EntityId      string
	EventRole     entities.EventRole
	ComponentType entities.ComponentType
}

var _ entities.Action = &Copy{}

func (c *Copy) Id() entities.ActionType {
	return entities.ActionCopy
}

func (c *Copy) Execute(ev *entities.Event) error {
	if ev.EntitiesById == nil {
		return fmt.Errorf("entities by id map in event may not be nil for copy action")
	}

	var recipient *entities.Entity
	switch c.EventRole {
	case entities.EventRoleSource:
		recipient = ev.Source
	case entities.EventRoleInstrument:
		recipient = ev.Instrument
	case entities.EventRoleTarget:
		recipient = ev.Target
	case entities.EventRoleRoom:
		recipient = ev.Room
	default:
		return fmt.Errorf("invalid role '%s' for copy action", c.EventRole.String())
	}

	component, err := recipient.RequireComponentWithChildren(c.ComponentType)
	if err != nil {
		return fmt.Errorf("error executing copy action: %w", err)
	}

	component.AddChild(
		ev.EntitiesById[c.EntityId].Copy(component),
	)

	return nil
}
