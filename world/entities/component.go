package entities

type Component interface {
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
