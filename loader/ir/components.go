// ir/components.go
package ir

import (
	"example.com/mud/world/entities"
	"example.com/mud/world/entities/components"
)

type Components struct {
	Room *RoomComponent
	// add more components over time...
}

func (c *Components) Build() map[string]entities.Component {
	out := make(map[string]entities.Component, 1)
	if c.Room != nil {
		out["room"] = c.Room.build()
	}
	return out
}

type RoomComponent struct {
	Color          string
	Icon           string
	Exits          map[string]string // dir -> room ID (these still resolve elsewhere)
	InlineChildren []*Entity         // <— inline child IRs (NOT IDs)
}

func (r *RoomComponent) build() *components.Room {
	return components.NewRoom(
		r.Icon,
		r.Color,
		cloneMap(r.Exits),
	)
}

func cloneMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
