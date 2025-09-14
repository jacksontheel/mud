package components

import (
	"strings"

	"example.com/mud/world/entities"
)

type Inventory struct {
	itemByAlias   map[string]*entities.Entity
	aliasesByItem map[*entities.Entity][]string
}

var _ entities.Component = &Inventory{}

func (i *Inventory) Copy() entities.Component {
	// right now copying inventories is not supported
	// return a new empty inventory
	return NewInventory([]*entities.Entity{})
}

func NewInventory(items []*entities.Entity) *Inventory {
	itemsByAlias := make(map[string]*entities.Entity, len(items)*2)
	aliasesByItem := make(map[*entities.Entity][]string, len(items))

	inventory := Inventory{
		itemByAlias:   itemsByAlias,
		aliasesByItem: aliasesByItem,
	}

	for _, item := range items {
		inventory.AddItem(item)
	}
	return &inventory
}

func (c *Inventory) GetItemByAlias(alias string) (*entities.Entity, bool) {
	entity, ok := c.itemByAlias[alias]
	return entity, ok
}

func (c *Inventory) AddItem(e *entities.Entity) {
	aliases := GetAliases(e)
	if len(aliases) == 0 {
		return
	}
	for _, alias := range aliases {
		c.aliasesByItem[e] = append(c.aliasesByItem[e], alias)
		c.itemByAlias[alias] = e
	}
}

func (c *Inventory) RemoveItem(e *entities.Entity) {
	aliases, ok := c.aliasesByItem[e]
	if !ok {
		return // abort if key does not exist
	}
	for _, alias := range aliases {
		delete(c.itemByAlias, alias) // for each alias, delete from itemsByAlias
	}
	delete(c.aliasesByItem, e) // delete entry from aliasesByItem
}

func (c *Inventory) Print() string {
	var b strings.Builder

	b.WriteString("You are carrying: [")

	for e := range c.aliasesByItem {
		if n := GetName(e); n != "" {
			b.WriteString(n)
			b.WriteString(", ")
		}
	}

	return strings.TrimSuffix(b.String(), ", ") + "]"
}
