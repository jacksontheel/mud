package components

import (
	"fmt"

	"example.com/mud/world/entities"
)

type Container struct {
	children *Children
}

var _ entities.Component = &Container{}
var _ entities.ComponentWithChildren = &Container{}

func (i *Container) Id() entities.ComponentType {
	return entities.ComponentContainer
}

func (i *Container) Copy() entities.Component {
	Container := NewContainer()
	for _, child := range i.children.GetChildren() {
		Container.children.AddChild(child)
	}
	return i
}

func NewContainer() *Container {
	return &Container{
		children: NewChildren(),
	}
}

func (c *Container) GetChildren() entities.IChildren {
	return c.children
}

func (c *Container) GetDescritionPrefix() string {
	return fmt.Sprintf("Inside the container:")
}
