package entities

type ComponentType int

const (
	ComponentRoom ComponentType = iota
	ComponentIdentity
	ComponentEventful
	ComponentInventory
)

type Component interface {
	Id() ComponentType
	Copy() Component
}

type ComponentWithChildren interface {
	GetChildren() IChildren
}

type IChildren interface {
	AddChild(child *Entity)
	RemoveChild(child *Entity)
	GetChildren() []*Entity
	GetChildByAlias(alias string) (*Entity, bool)
}
