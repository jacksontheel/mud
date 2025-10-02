package components

import (
	"fmt"
	"strings"

	"example.com/mud/world/entities"
)

type Inventory struct {
	children *Children
}

var _ entities.Component = &Inventory{}
var _ entities.ComponentWithChildren = &Inventory{}

func (i *Inventory) Id() entities.ComponentType {
	return entities.ComponentInventory
}

func (i *Inventory) Copy() entities.Component {
	inventory := NewInventory()
	for _, child := range i.children.GetChildren() {
		inventory.children.AddChild(child)
	}
	return i
}

func NewInventory() *Inventory {
	return &Inventory{
		children: NewChildren(),
	}
}

func (i *Inventory) GetChildren() entities.IChildren {
	return i.children
}

func (i *Inventory) GetDescritionPrefix() string {
	return fmt.Sprintf("Inside the inventory:")
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
