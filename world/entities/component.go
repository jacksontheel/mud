package entities

import "fmt"

type ComponentType int

const (
	ComponentUnknown ComponentType = iota
	ComponentRoom
	ComponentEventful
	ComponentInventory
)

const (
	ComponentUnknownString   = "Unknown"
	ComponentRoomString      = "Room"
	ComponentEventfulString  = "Eventful"
	ComponentInventoryString = "Inventory"
)

func ParseComponentType(s string) (ComponentType, error) {
	switch s {
	case ComponentRoomString:
		return ComponentRoom, nil
	case ComponentEventfulString:
		return ComponentEventful, nil
	case ComponentInventoryString:
		return ComponentInventory, nil
	default:
		return ComponentUnknown, fmt.Errorf("unknown component type '%s'", s)
	}
}

func (ct ComponentType) String() string {
	switch ct {
	case ComponentRoom:
		return ComponentRoomString
	case ComponentEventful:
		return ComponentEventfulString
	case ComponentInventory:
		return ComponentInventoryString
	default:
		return ComponentUnknownString
	}
}

type Component interface {
	Id() ComponentType
	Copy() Component
}

type ComponentWithChildren interface {
	GetChildren() IChildren
}

type IChildren interface {
	AddChild(child *Entity) error
	RemoveChild(child *Entity)
	GetChildren() []*Entity
	GetChildByAlias(alias string) (*Entity, bool)
	HasChild(e *Entity) bool
}
