package actions

import (
	"fmt"

	"example.com/mud/world/entities"
)

type Move struct {
	RoleOrigin      EventRole
	RoleDestination EventRole
	ComponentType   entities.ComponentType
}

var _ entities.Action = &Move{}

func (m *Move) Id() entities.ActionType {
	return entities.ActionMove
}

func (m *Move) Execute(ev *entities.Event) error {

	var origin *entities.Entity
	switch m.RoleOrigin {
	case EventRoleSource:
		origin = ev.Source
	case EventRoleInstrument:
		origin = ev.Instrument
	case EventRoleTarget:
		origin = ev.Target
	default:
		return fmt.Errorf("invalid origin role '%s' for move action", m.RoleOrigin.String())
	}

	if origin == nil {
		return fmt.Errorf("origin role '%s' is empty for move event", m.RoleOrigin)
	}

	var destination *entities.Entity
	switch m.RoleDestination {
	case EventRoleSource:
		destination = ev.Source
	case EventRoleInstrument:
		destination = ev.Instrument
	case EventRoleTarget:
		destination = ev.Target
	default:
		return fmt.Errorf("invalid origin role '%s' for move action", m.RoleDestination.String())
	}

	if destination == nil {
		return fmt.Errorf("destination role '%s' is empty for move event", m.RoleOrigin)
	}

	component, err := destination.RequireComponentWithChildren(m.ComponentType)
	if err != nil {
		return fmt.Errorf("error executing copy action: %w", err)
	}

	// remove entity from old parent
	oldParent := origin.Parent
	oldParent.GetChildren().RemoveChild(origin)

	// set parent of entity and add it to component
	origin.Parent = component
	component.GetChildren().AddChild(origin)

	return nil
}
