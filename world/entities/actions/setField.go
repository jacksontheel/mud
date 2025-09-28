package actions

import (
	"fmt"

	"example.com/mud/world/entities"
)

type SetField struct {
	Role  entities.EventRole
	Field string
	Value any
}

var _ entities.Action = &SetField{}

func (sf *SetField) Id() entities.ActionType {
	return entities.ActionSetField
}

func (sf *SetField) Execute(ev *entities.Event) error {

	var e *entities.Entity
	switch sf.Role {
	case entities.EventRoleSource:
		e = ev.Source
	case entities.EventRoleInstrument:
		e = ev.Instrument
	case entities.EventRoleTarget:
		e = ev.Target
	default:
		return fmt.Errorf("invalid role '%s' for SetField action", sf.Role)
	}

	if e == nil {
		return fmt.Errorf("role '%s' is empty for SetField event", sf.Role)
	}

	e.Fields[sf.Field] = sf.Value

	return nil
}
