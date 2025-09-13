package entities

import (
	"strings"
)

type Inventory struct {
	itemByAlias   map[string]*Entity
	aliasesByItem map[*Entity][]string
}

var _ Component = &Inventory{}

func NewInventory(items []*Entity) *Inventory {
	itemsByAlias := make(map[string]*Entity, len(items)*2)
	aliasesByItem := make(map[*Entity][]string, len(items))

	inventory := Inventory{
		itemByAlias:   itemsByAlias,
		aliasesByItem: aliasesByItem,
	}

	for _, item := range items {
		inventory.AddItem(item)
	}
	return &inventory
}

func (i *Inventory) Id() string {
	return "inventory"
}

func (c *Inventory) GetItemByAlias(alias string) (*Entity, bool) {
	entity, ok := c.itemByAlias[alias]
	return entity, ok
}

func (c *Inventory) AddItem(e *Entity) {
	aliases := e.getAliases()
	if len(aliases) == 0 {
		return
	}
	for _, alias := range aliases {
		key := normalizeAlias(alias)
		if key == "" {
			continue
		}
		c.aliasesByItem[e] = append(c.aliasesByItem[e], key)
		c.itemByAlias[key] = e
	}
}

func (c *Inventory) RemoveItem(e *Entity) {
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
		if n := e.getName(); n != "" {
			b.WriteString(n)
			b.WriteString(", ")
		}
	}

	return strings.TrimSuffix(b.String(), ", ") + "]"
}
