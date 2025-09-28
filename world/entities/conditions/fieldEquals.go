package conditions

import (
	"fmt"

	"example.com/mud/world/entities"
)

type FieldEquals struct {
	Role  entities.EventRole
	Field string
	Value any
}

var _ entities.Condition = &FieldEquals{}

func (fe *FieldEquals) Id() entities.ConditionType {
	return entities.ConditionFieldEquals
}

func (fe *FieldEquals) Check(ev *entities.Event) (bool, error) {
	var e *entities.Entity
	switch fe.Role {
	case entities.EventRoleSource:
		e = ev.Source
	case entities.EventRoleInstrument:
		e = ev.Instrument
	case entities.EventRoleTarget:
		e = ev.Target
	default:
		return false, fmt.Errorf("invalid role '%s' for field equals condition", fe.Role.String())
	}

	value, ok := e.Fields[fe.Field]
	if !ok {
		return false, fmt.Errorf("entity '%s' does not have field '%s'", e.Name, fe.Field)
	}

	return value == fe.Value, nil
}
