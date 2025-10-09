package components

import "example.com/mud/world/entities"

type Children struct {
	revealed bool
	prefix   string

	childByAlias   map[string]*entities.Entity
	aliasesByChild map[*entities.Entity][]string
}

var _ entities.IChildren = &Children{}

func NewChildren() *Children {
	return &Children{
		childByAlias:   make(map[string]*entities.Entity),
		aliasesByChild: make(map[*entities.Entity][]string),
	}
}

func (c *Children) Copy() entities.IChildren {
	copiedChildren := NewChildren()

	copiedChildren.revealed = c.revealed
	copiedChildren.prefix = c.prefix

	for _, child := range c.GetChildren() {
		copiedChildren.AddChild(child)
	}

	return copiedChildren
}

func (c *Children) GetPrefix() string {
	return c.prefix
}

func (c *Children) GetRevealed() bool {
	return c.revealed
}

func (c *Children) SetPrefix(p string) {
	c.prefix = p
}

func (c *Children) SetRevealed(r bool) {
	c.revealed = r
}

func (c *Children) AddChild(child *entities.Entity) error {
	aliases := child.Aliases

	if len(aliases) == 0 {
		return nil
	}
	for _, alias := range aliases {
		c.aliasesByChild[child] = append(c.aliasesByChild[child], alias)
		c.childByAlias[alias] = child
	}

	return nil
}

func (c *Children) RemoveChild(child *entities.Entity) {
	aliases, ok := c.aliasesByChild[child]
	if !ok {
		return
	}

	for _, alias := range aliases {
		delete(c.childByAlias, alias) // for each alias, delete from itemsByAlias
	}
	delete(c.aliasesByChild, child) // delete entry from aliasesByItem
}

func (c *Children) GetChildren() []*entities.Entity {
	children := make([]*entities.Entity, 0)
	for child := range c.aliasesByChild {
		children = append(children, child)
	}
	return children
}

func (c *Children) GetChildByAlias(alias string) (*entities.Entity, bool) {
	child, ok := c.childByAlias[alias]
	if ok {
		return child, ok
	}

	for _, children := range c.GetChildren() {
		for _, cwc := range children.GetComponentsWithChildren() {
			if !cwc.GetChildren().GetRevealed() {
				continue
			}

			grandchild, ok := cwc.GetChildren().GetChildByAlias(alias)
			return grandchild, ok
		}
	}

	return child, ok
}

func (c *Children) HasChild(e *entities.Entity) bool {
	_, ok := c.aliasesByChild[e]

	return ok
}
