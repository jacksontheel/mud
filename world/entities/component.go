package entities

type ComponentType int

const (
	ComponentUnknown ComponentType = iota
	ComponentRoom
	ComponentIdentity
	ComponentEventful
	ComponentInventory
)

const (
	ComponentUnknownString   = "Unknown"
	ComponentRoomString      = "Room"
	ComponentIdentityString  = "Identity"
	ComponentEventfulString  = "Eventful"
	ComponentInventoryString = "Inventory"
)

func (ct ComponentType) String() string {
	switch ct {
	case ComponentRoom:
		return ComponentRoomString
	case ComponentIdentity:
		return ComponentIdentityString
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
}
