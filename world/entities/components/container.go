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

func NewContainer() *Container {
	return &Container{
		children: NewChildren(),
	}
}

func (i *Container) Id() entities.ComponentType {
	return entities.ComponentContainer
}

func (i *Container) Copy() entities.Component {
	Container := NewContainer()
	for _, child := range i.children.GetChildren() {
		Container.AddChild(child)
	}
	return i
}

func (c *Container) AddChild(child *entities.Entity) error {
	err := c.GetChildren().AddChild(child)
	if err != nil {
		return fmt.Errorf("Inventory add child: %w", err)
	}

	child.Parent = c

	return nil
}

func (c *Container) RemoveChild(child *entities.Entity) {
	child.Parent = nil
	c.GetChildren().RemoveChild(child)
}

func (c *Container) GetChildren() entities.IChildren {
	return c.children
}
