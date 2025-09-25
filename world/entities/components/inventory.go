package components

import (
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
	// right now copying inventories is not supported
	// return a new empty inventory
	return NewInventory()
}

func NewInventory() *Inventory {
	return &Inventory{
		children: NewChildren(),
	}
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
