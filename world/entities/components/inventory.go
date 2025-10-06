package components

import (
	"fmt"
	"strings"

	"example.com/mud/world/entities"
)

type Inventory struct {
	children *Children
}

func NewInventory() *Inventory {
	return &Inventory{
		children: NewChildren(),
	}
}

var _ entities.Component = &Inventory{}
var _ entities.ComponentWithChildren = &Inventory{}

func (i *Inventory) Id() entities.ComponentType {
	return entities.ComponentInventory
}

func (i *Inventory) Copy() entities.Component {
	inventory := NewInventory()
	for _, child := range i.children.GetChildren() {
		inventory.AddChild(child)
	}
	return inventory
}

func (i *Inventory) AddChild(child *entities.Entity) error {
	err := i.GetChildren().AddChild(child)
	if err != nil {
		return fmt.Errorf("Inventory add child: %w", err)
	}

	child.Parent = i

	return nil
}

func (i *Inventory) RemoveChild(child *entities.Entity) {
	child.Parent = nil
	i.GetChildren().RemoveChild(child)
}

func (i *Inventory) GetChildren() entities.IChildren {
	return i.children
}

func (i *Inventory) Print() (string, error) {
	var b strings.Builder

	b.WriteString("You are carrying: [")

	for _, child := range i.GetChildren().GetChildren() {
		if n := child.Name; n != "" {
			b.WriteString(n)
			b.WriteString(", ")
		}
	}

	return strings.TrimSuffix(b.String(), ", ") + "]", nil
}
