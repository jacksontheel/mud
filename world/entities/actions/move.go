package actions

import (
	"fmt"

	"example.com/mud/world/entities"
)

type Move struct {
	RoleOrigin      entities.EventRole
	RoleDestination entities.EventRole
	ComponentType   entities.ComponentType
}

var _ entities.Action = &Move{}

func (m *Move) Id() entities.ActionType {
	return entities.ActionMove
}

func (m *Move) Execute(ev *entities.Event) error {

	var origin *entities.Entity
	switch m.RoleOrigin {
	case entities.EventRoleSource:
		origin = ev.Source
	case entities.EventRoleInstrument:
		origin = ev.Instrument
	case entities.EventRoleTarget:
		origin = ev.Target
	case entities.EventRoleRoom:
		origin = ev.Room
	default:
		return fmt.Errorf("invalid origin role '%s' for move action", m.RoleOrigin.String())
	}

	if origin == nil {
		return fmt.Errorf("origin role '%s' is empty for move event", m.RoleOrigin)
	}

	var destination *entities.Entity
	switch m.RoleDestination {
	case entities.EventRoleSource:
		destination = ev.Source
	case entities.EventRoleInstrument:
		destination = ev.Instrument
	case entities.EventRoleTarget:
		destination = ev.Target
	case entities.EventRoleRoom:
		destination = ev.Room
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
	oldParent.RemoveChild(origin)

	// add entity to new parent
	component.AddChild(origin)

	return nil
}
